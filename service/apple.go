package service

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"microservices/entity/config"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	appleVerifyReceiptProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
	appleVerifyReceiptSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	appleRootCAPath                 = "AppleRootCA-G3.cer"
)

var (
	appleRootCACert *x509.Certificate
	rootCALoadOnce  sync.Once
	rootCALoadError error
)

// Apple exposes integrations with Apple's in-app purchase APIs.
type Apple interface {
	VerifyReceipt(ctx context.Context, receiptData string, excludeOldTransactions bool) (*AppleVerifyReceiptResult, error)
	DecodeSignedPayload(signedPayload string) (*AppleNotificationPayload, error)
	DecodeSignedTransaction(signedTransaction string) (*AppleTransactionInfo, error)
	DecodeSignedRenewalInfo(signedRenewalInfo string) (*AppleRenewalInfo, error)
}

type appleService struct {
	httpClient *http.Client
	options    *config.AppleOptions
}

type appleVerifyReceiptRequest struct {
	ReceiptData            string `json:"receipt-data"`
	Password               string `json:"password"`
	ExcludeOldTransactions bool   `json:"exclude-old-transactions,omitempty"`
}

type appleVerifyReceiptResponse struct {
	Status            int                   `json:"status"`
	Environment       string                `json:"environment"`
	LatestReceipt     string                `json:"latest_receipt"`
	LatestReceiptInfo []appleReceiptInfoRaw `json:"latest_receipt_info"`
}

type appleReceiptInfoRaw struct {
	ProductID             string `json:"product_id"`
	TransactionID         string `json:"transaction_id"`
	OriginalTransactionID string `json:"original_transaction_id"`
	PurchaseDateMs        string `json:"purchase_date_ms"`
	ExpiresDateMs         string `json:"expires_date_ms"`
	CancellationDateMs    string `json:"cancellation_date_ms"`
}

// AppleReceiptInfo normalised receipt information returned from Apple.
type AppleReceiptInfo struct {
	ProductID             string
	TransactionID         string
	OriginalTransactionID string
	PurchaseDateMs        string
	ExpiresDateMs         string
	CancellationDateMs    string
}

// AppleVerifyReceiptResult wraps the relevant verify receipt response values.
type AppleVerifyReceiptResult struct {
	Status            int
	Environment       string
	LatestReceipt     string
	LatestReceiptInfo []AppleReceiptInfo
}

func newApple(opt *config.AppleOptions) Apple {
	return &appleService{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		options:    opt,
	}
}

func (a *appleService) VerifyReceipt(ctx context.Context, receiptData string, excludeOldTransactions bool) (*AppleVerifyReceiptResult, error) {
	if receiptData == "" {
		return nil, fmt.Errorf("receipt data is empty")
	}
	if a.options == nil || a.options.SharedSecret == "" {
		return nil, fmt.Errorf("apple shared secret is not configured")
	}

	reqPayload := &appleVerifyReceiptRequest{
		ReceiptData: receiptData,
		Password:    a.options.SharedSecret,
	}
	if excludeOldTransactions {
		reqPayload.ExcludeOldTransactions = true
	}

	result, err := a.sendVerifyReceipt(ctx, appleVerifyReceiptProductionURL, reqPayload)
	if err != nil {
		return nil, err
	}

	// If status 21007 returned, retry against sandbox endpoint.
	if result.Status == 21007 {
		result, err = a.sendVerifyReceipt(ctx, appleVerifyReceiptSandboxURL, reqPayload)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (a *appleService) sendVerifyReceipt(ctx context.Context, endpoint string, payload *appleVerifyReceiptRequest) (*AppleVerifyReceiptResult, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal verify receipt request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create verify receipt request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call verify receipt: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("verify receipt http status %d: %s", resp.StatusCode, bytes.TrimSpace(data))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read verify receipt response: %w", err)
	}

	var parsed appleVerifyReceiptResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("decode verify receipt response: %w", err)
	}

	result := &AppleVerifyReceiptResult{
		Status:        parsed.Status,
		Environment:   parsed.Environment,
		LatestReceipt: parsed.LatestReceipt,
	}

	if len(parsed.LatestReceiptInfo) > 0 {
		result.LatestReceiptInfo = make([]AppleReceiptInfo, 0, len(parsed.LatestReceiptInfo))
		for _, info := range parsed.LatestReceiptInfo {
			result.LatestReceiptInfo = append(result.LatestReceiptInfo, AppleReceiptInfo{
				ProductID:             info.ProductID,
				TransactionID:         info.TransactionID,
				OriginalTransactionID: info.OriginalTransactionID,
				PurchaseDateMs:        info.PurchaseDateMs,
				ExpiresDateMs:         info.ExpiresDateMs,
				CancellationDateMs:    info.CancellationDateMs,
			})
		}
	}

	return result, nil
}

// ==================== JWS (JSON Web Signature) Support for App Store Server Notifications V2 ====================

// jwsHeader represents the header of a JWS token
type jwsHeader struct {
	Alg string   `json:"alg"`
	X5c []string `json:"x5c"` // Certificate chain
}

// AppleNotificationPayload represents the decoded notification data from Apple
type AppleNotificationPayload struct {
	NotificationType string          `json:"notificationType"`
	Subtype          string          `json:"subtype"`
	NotificationUUID string          `json:"notificationUUID"`
	Data             json.RawMessage `json:"data"`
	Version          string          `json:"version"`
	SignedDate       int64           `json:"signedDate"`
}

// AppleTransactionInfo represents the signed transaction information
type AppleTransactionInfo struct {
	AppAccountToken             string `json:"appAccountToken"`
	TransactionId               string `json:"transactionId"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	WebOrderLineItemId          string `json:"webOrderLineItemId"`
	BundleId                    string `json:"bundleId"`
	ProductId                   string `json:"productId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	PurchaseDate                int64  `json:"purchaseDate"`
	OriginalPurchaseDate        int64  `json:"originalPurchaseDate"`
	ExpiresDate                 int64  `json:"expiresDate"`
	Quantity                    int    `json:"quantity"`
	Type                        string `json:"type"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	SignedDate                  int64  `json:"signedDate"`
	Environment                 string `json:"environment"`
}

// AppleRenewalInfo represents the signed renewal information
type AppleRenewalInfo struct {
	ExpirationIntent            int    `json:"expirationIntent"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	AutoRenewProductId          string `json:"autoRenewProductId"`
	ProductId                   string `json:"productId"`
	AutoRenewStatus             int    `json:"autoRenewStatus"`
	IsInBillingRetryPeriod      bool   `json:"isInBillingRetryPeriod"`
	SignedDate                  int64  `json:"signedDate"`
	Environment                 string `json:"environment"`
	RecentSubscriptionStartDate int64  `json:"recentSubscriptionStartDate"`
}

// DecodeSignedPayload decodes and verifies the JWS signed payload from Apple
// It extracts the certificate chain from the header and verifies the signature
func (a *appleService) DecodeSignedPayload(signedPayload string) (*AppleNotificationPayload, error) {
	publicKey, err := a.extractPublicKeyFromJWS(signedPayload)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(signedPayload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify JWS signature: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWS token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	var payload AppleNotificationPayload
	if err := json.Unmarshal(claimsJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}
	if time.UnixMilli(payload.SignedDate).Add(5 * time.Minute).Before(time.Now()) {
		return nil, fmt.Errorf("signedDate expired")
	}

	return &payload, nil
}

// DecodeSignedTransaction decodes the signed transaction info
func (a *appleService) DecodeSignedTransaction(signedTransaction string) (*AppleTransactionInfo, error) {
	publicKey, err := a.extractPublicKeyFromJWS(signedTransaction)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(signedTransaction, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify JWS signature: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWS token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	var transactionInfo AppleTransactionInfo
	if err := json.Unmarshal(claimsJSON, &transactionInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction info: %w", err)
	}

	return &transactionInfo, nil
}

// DecodeSignedRenewalInfo decodes the signed renewal info
func (a *appleService) DecodeSignedRenewalInfo(signedRenewalInfo string) (*AppleRenewalInfo, error) {
	publicKey, err := a.extractPublicKeyFromJWS(signedRenewalInfo)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(signedRenewalInfo, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify JWS signature: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWS token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claims: %w", err)
	}

	var renewalInfo AppleRenewalInfo
	if err := json.Unmarshal(claimsJSON, &renewalInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal renewal info: %w", err)
	}

	return &renewalInfo, nil
}

// extractPublicKeyFromJWS extracts and verifies the ECDSA public key from the JWS header's certificate chain
func (a *appleService) extractPublicKeyFromJWS(jws string) (*ecdsa.PublicKey, error) {
	parts := strings.Split(jws, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format: expected 3 parts, got %d", len(parts))
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWS header: %w", err)
	}

	var header jwsHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWS header: %w", err)
	}

	if len(header.X5c) == 0 {
		return nil, fmt.Errorf("no certificate chain found in JWS header")
	}

	// Parse the complete certificate chain from X5c
	certChain := make([]*x509.Certificate, 0, len(header.X5c))
	for i, certBase64 := range header.X5c {
		certBytes, err := base64.StdEncoding.DecodeString(certBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode certificate %d: %w", i, err)
		}

		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate %d: %w", i, err)
		}

		certChain = append(certChain, cert)
	}

	// Verify the certificate chain against Apple Root CA
	if err := verifyCertificateChain(certChain); err != nil {
		return nil, fmt.Errorf("certificate chain verification failed: %w", err)
	}

	// Extract the public key from the leaf certificate (first in chain)
	publicKey, ok := certChain[0].PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate public key is not ECDSA")
	}

	return publicKey, nil
}

// ==================== Certificate Verification ====================

// loadAppleRootCA loads the Apple Root CA certificate from file (once)
func loadAppleRootCA() (*x509.Certificate, error) {
	var cert *x509.Certificate
	var loadErr error

	rootCALoadOnce.Do(func() {
		certData, err := os.ReadFile(appleRootCAPath)
		if err != nil {
			loadErr = fmt.Errorf("failed to read Apple Root CA file: %w", err)
			return
		}

		parsedCert, err := x509.ParseCertificate(certData)
		if err != nil {
			loadErr = fmt.Errorf("failed to parse Apple Root CA certificate: %w", err)
			return
		}

		appleRootCACert = parsedCert
		cert = parsedCert
	})

	if loadErr != nil {
		return nil, loadErr
	}

	if appleRootCACert != nil {
		return appleRootCACert, nil
	}

	return cert, rootCALoadError
}

// verifyCertificateChain verifies the complete certificate chain against Apple Root CA
func verifyCertificateChain(certChain []*x509.Certificate) error {
	if len(certChain) == 0 {
		return fmt.Errorf("empty certificate chain")
	}

	// Load Apple Root CA certificate
	rootCA, err := loadAppleRootCA()
	if err != nil {
		return fmt.Errorf("failed to load Apple Root CA: %w", err)
	}

	// Create a certificate pool with Apple Root CA
	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootCA)

	// Create intermediate pool with all certs except the leaf (first) certificate
	intermediatePool := x509.NewCertPool()
	for i := 1; i < len(certChain); i++ {
		intermediatePool.AddCert(certChain[i])
	}

	// Verify the leaf certificate
	leafCert := certChain[0]
	verifyOpts := x509.VerifyOptions{
		Roots:         rootPool,
		Intermediates: intermediatePool,
		CurrentTime:   time.Now(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	chains, err := leafCert.Verify(verifyOpts)
	if err != nil {
		return fmt.Errorf("certificate verification failed: %w", err)
	}

	if len(chains) == 0 {
		return fmt.Errorf("no valid certificate chains found")
	}

	// Additional validation: check certificate is from Apple
	if len(leafCert.Subject.Organization) == 0 {
		return fmt.Errorf("certificate has no organization")
	}

	// Check that the certificate is issued by Apple
	foundApple := false
	for _, org := range leafCert.Subject.Organization {
		if strings.Contains(org, "Apple") {
			foundApple = true
			break
		}
	}

	if !foundApple {
		return fmt.Errorf("certificate is not issued by Apple")
	}

	return nil
}
