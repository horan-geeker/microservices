package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type System interface {
	CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int64, error)
	GetRateLimit(ctx context.Context, key string) (int64, error)
	HandleLoginFailure(ctx context.Context, userID int) (bool, error)
	IsAccountFrozen(ctx context.Context, userID int) (bool, error)
	ClearLoginFailure(ctx context.Context, userID int) error
}

type system struct {
	rdb *redis.Client
}

func (s *system) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int64, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	// 使用 INCR 增加计数
	count, err := s.rdb.Incr(ctx, rateLimitKey).Result()
	if err != nil {
		return false, 0, err
	}

	// 如果是第一次请求，设置过期时间
	if count == 1 {
		err = s.rdb.Expire(ctx, rateLimitKey, window).Err()
		if err != nil {
			return false, count, err
		}
	}

	// 检查是否超过限制
	if count > int64(limit) {
		return false, count, nil
	}

	return true, count, nil
}

func (s *system) GetRateLimit(ctx context.Context, key string) (int64, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)
	result, err := s.rdb.Get(ctx, rateLimitKey).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	count, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// HandleLoginFailure 处理登录失败，达到5次自动冻结2小时，返回是否已冻结
func (s *system) HandleLoginFailure(ctx context.Context, userID int) (bool, error) {
	failureKey := fmt.Sprintf("login_failure:%d", userID)
	frozenKey := fmt.Sprintf("account_frozen:%d", userID)

	// 增加失败次数
	count, err := s.rdb.Incr(ctx, failureKey).Result()
	if err != nil {
		return false, err
	}

	// 如果是第一次失败，设置24小时过期（失败记录保留时间）
	if count == 1 {
		err = s.rdb.Expire(ctx, failureKey, 24*time.Hour).Err()
		if err != nil {
			return false, err
		}
	}

	// 如果达到5次失败，冻结账号2小时
	if count >= 5 {
		err = s.rdb.Set(ctx, frozenKey, "1", 2*time.Hour).Err()
		if err != nil {
			return false, err
		}
		// 清除失败计数
		s.rdb.Del(ctx, failureKey)
		return true, nil
	}

	return false, nil
}

// IsAccountFrozen 检查账号是否被冻结
func (s *system) IsAccountFrozen(ctx context.Context, userID int) (bool, error) {
	frozenKey := fmt.Sprintf("account_frozen:%d", userID)

	result, err := s.rdb.Get(ctx, frozenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return result == "1", nil
}

// ClearLoginFailure 清除登录失败计数
func (s *system) ClearLoginFailure(ctx context.Context, userID int) error {
	failureKey := fmt.Sprintf("login_failure:%d", userID)
	return s.rdb.Del(ctx, failureKey).Err()
}

func newSystem(rdb *redis.Client) System {
	return &system{rdb: rdb}
}
