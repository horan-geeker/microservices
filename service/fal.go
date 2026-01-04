package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservices/entity/config"
	"net/http"
	"os"
	"time"
)

type FalService interface {
	RemoveBackground(ctx context.Context, imageUrl string) (string, error)
}

type falService struct {
	config *config.FalOptions
}

func NewFalService(opt *config.FalOptions) FalService {
	return &falService{
		config: opt,
	}
}

type falQueueResponse struct {
	RequestID string `json:"request_id"`
}

type falStatusResponse struct {
	Status  string `json:"status"` // IN_QUEUE, IN_PROGRESS, COMPLETED, FAILED
	Logs    []any  `json:"logs"`
	Metrics any    `json:"metrics"`
	Image   struct {
		Url         string `json:"url"`
		ContentType string `json:"content_type"`
		FileName    string `json:"file_name"`
		FileSize    int    `json:"file_size"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
	} `json:"image"`
	Error string `json:"error"`
}

func (s *falService) RemoveBackground(ctx context.Context, imageUrl string) (string, error) {
	apiKey := os.Getenv("FAL_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("FAL_KEY not set")
	}

	url := "https://queue.fal.run/fal-ai/bria/background/remove"
	payload := map[string]string{
		"image_url": imageUrl,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("fal api error: %s", string(body))
	}

	var queueResp falQueueResponse
	if err := json.NewDecoder(resp.Body).Decode(&queueResp); err != nil {
		return "", err
	}

	if queueResp.RequestID == "" {
		return "", fmt.Errorf("no request_id returned from fal")
	}

	// Poll for result
	return s.pollResult(ctx, apiKey, queueResp.RequestID)
}

func (s *falService) pollResult(ctx context.Context, apiKey, requestID string) (string, error) {
	url := fmt.Sprintf("https://queue.fal.run/fal-ai/bria/background/remove/requests/%s", requestID)
	client := &http.Client{}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timeout:
			return "", fmt.Errorf("fal processing timeout")
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return "", err
			}
			req.Header.Set("Authorization", "Key "+apiKey)

			resp, err := client.Do(req)
			if err != nil {
				// Log error but continue polling? or fail? Let's fail for now.
				return "", err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				body, _ := io.ReadAll(resp.Body)
				return "", fmt.Errorf("fal poll error: %s", string(body))
			}

			// Read body to string for debugging/parsing flexibility since some fields might change
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}

			var statusResp falStatusResponse
			if err := json.Unmarshal(bodyBytes, &statusResp); err != nil {
				return "", err
			}

			if statusResp.Status == "COMPLETED" {
				if statusResp.Image.Url != "" {
					return statusResp.Image.Url, nil
				}
				return "", fmt.Errorf("completed but no image url found")
			} else if statusResp.Status == "FAILED" {
				return "", fmt.Errorf("fal processing failed: %s", statusResp.Error)
			}
			// Continue polling if IN_QUEUE or IN_PROGRESS
		}
	}
}
