package redis

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const searchCacheTTL = 3600 * time.Second

type SearchCache struct {
	client *redis.Client
}

func NewSearchCache(client *redis.Client) *SearchCache {
	return &SearchCache{client: client}
}

func searchCacheKey(query string) string {
	return fmt.Sprintf("%s%x", cachePrefix, sha256.Sum256([]byte(query)))
}

func (c *SearchCache) Get(ctx context.Context, query string) ([]byte, bool, error) {
	body, err := c.client.Get(ctx, searchCacheKey(query)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("get search cache: %w", err)
	}
	return body, true, nil
}

func (c *SearchCache) Set(ctx context.Context, query string, body []byte) error {
	if err := c.client.Set(ctx, searchCacheKey(query), body, searchCacheTTL).Err(); err != nil {
		return fmt.Errorf("set search cache: %w", err)
	}
	return nil
}
