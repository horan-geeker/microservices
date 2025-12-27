package callback

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"microservices/entity/config"
	"microservices/entity/consts"
	"microservices/entity/jwt"
	entity "microservices/entity/model"
	"time"
)

func (l *logic) GoogleCallback(ctx context.Context, code string) (*entity.User, string, error) {
	// 1. Exchange code for token
	token, err := l.srv.Google().ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("exchange code failed: %w", err)
	}

	// 2. Get User Info
	tokenSource := oauth2.StaticTokenSource(token)
	userInfo, err := l.srv.Google().GetUserInfo(ctx, tokenSource)
	if err != nil {
		return nil, "", fmt.Errorf("get user info failed: %w", err)
	}

	// 3. Check Authorize
	provider := "google"
	providerID := userInfo.Id
	auth, err := l.model.Authorize().GetByProvider(ctx, provider, providerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	var user *entity.User

	if auth != nil && auth.ID != 0 {
		// Found authorized record, get user
		user, err = l.model.User().GetByUid(ctx, auth.UserID)
		if err != nil {
			return nil, "", err
		}
		// Update auth token info
		authData := map[string]any{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
			"expires_at":    token.Expiry,
			"updated_at":    time.Now(),
		}
		jsonBytes, _ := json.Marshal(userInfo)
		authData["json"] = string(jsonBytes)
		_ = l.model.Authorize().Update(ctx, auth.ID, authData)

		// Sync User Avatar if changed or empty
		if userInfo.Picture != "" && user.Avatar != userInfo.Picture {
			_ = l.model.User().Update(ctx, user.ID, map[string]any{
				"avatar": userInfo.Picture,
			})
			user.Avatar = userInfo.Picture // update local object
		}
	} else {
		// Not found, check existing user by email
		user, err = l.model.User().GetByEmail(ctx, userInfo.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", err
		}

		if user == nil || user.ID == 0 {
			// Create new user

			newUser := &entity.User{
				Name:      userInfo.Name,
				Avatar:    userInfo.Picture,
				Email:     userInfo.Email,
				Password:  "",
				Status:    consts.UserStatusNormal,
				LoginAt:   time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := l.model.User().Create(ctx, newUser); err != nil {
				return nil, "", err
			}
			user = newUser
		} else {
			// Exist user found by email, sync avatar if needed
			if userInfo.Picture != "" && user.Avatar != userInfo.Picture {
				_ = l.model.User().Update(ctx, user.ID, map[string]any{
					"avatar": userInfo.Picture,
				})
				user.Avatar = userInfo.Picture
			}
		}

		// Create Authorize record
		jsonBytes, _ := json.Marshal(userInfo)
		newAuth := &entity.Authorize{
			UserID:       user.ID,
			Provider:     provider,
			ProviderID:   providerID,
			Email:        userInfo.Email,
			Name:         userInfo.Name,
			Avatar:       userInfo.Picture,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.Expiry,
			JSON:         string(jsonBytes),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := l.model.Authorize().Create(ctx, newAuth); err != nil {
			return nil, "", err
		}
	}

	// 4. Login (Generate Token)
	jwtToken, err := jwt.NewJwt(config.NewJwtOptions()).GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	if err := l.cache.User().SetToken(ctx, user.ID, jwtToken); err != nil {
		return nil, "", err
	}
	// Update login time
	_ = l.model.User().Update(ctx, user.ID, map[string]any{
		"login_at": time.Now(),
	})

	return user, jwtToken, nil
}
