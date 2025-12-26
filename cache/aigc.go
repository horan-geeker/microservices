package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AIGC interface {
	SetGenerationLock(ctx context.Context, uid int) (bool, error)
	ReleaseGenerationLock(ctx context.Context, uid int) error
}

type aigc struct {
	rdb *redis.Client
}

func newAIGC(rdb *redis.Client) AIGC {
	return &aigc{rdb: rdb}
}

func (c *aigc) SetGenerationLock(ctx context.Context, uid int) (bool, error) {
	key := fmt.Sprintf("lock:aigc:generation:%d", uid)
	// Lock expires in 5 minutes to prevent deadlocks if something crashes hard
	return c.rdb.SetNX(ctx, key, 1, 5*time.Minute).Result()
}

func (c *aigc) ReleaseGenerationLock(ctx context.Context, uid int) error {
	key := fmt.Sprintf("lock:aigc:generation:%d", uid)
	return c.rdb.Del(ctx, key).Err()
}
