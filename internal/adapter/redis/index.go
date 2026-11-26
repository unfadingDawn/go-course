package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const indexPrefix = "movies:index:"
const cachePrefix = "movies:cache:"

type MovieIndex struct {
	client *redis.Client
}

func NewMovieIndex(client *redis.Client) *MovieIndex {
	return &MovieIndex{client: client}
}

func (i *MovieIndex) AddMovie(ctx context.Context, id int, words []string) error {
	if len(words) == 0 {
		return nil
	}

	_, err := i.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, word := range words {
			key := indexPrefix + word
			pipe.SAdd(ctx, key, id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("index movie %d: %w", id, err)
	}
	return nil
}

func (i *MovieIndex) Reset(ctx context.Context) error {
	for _, pattern := range []string{indexPrefix + "*", cachePrefix + "*"} {
		var cursor uint64
		for {
			keys, next, err := i.client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return fmt.Errorf("scan keys %s: %w", pattern, err)
			}
			for start := 0; start < len(keys); start += 100 {
				end := min(start+100, len(keys))
				if err := i.client.Del(ctx, keys[start:end]...).Err(); err != nil {
					return fmt.Errorf("delete keys %s: %w", pattern, err)
				}
			}
			cursor = next
			if cursor == 0 {
				break
			}
		}
	}
	return nil
}

func (i *MovieIndex) FindIDs(ctx context.Context, words []string) ([]int, error) {
	if len(words) == 0 {
		return []int{}, nil
	}
	keys := make([]string, 0, len(words))
	for _, word := range words {
		keys = append(keys, indexPrefix+word)
	}
	members, err := i.client.SInter(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("find movie IDs: %w", err)
	}
	ids := make([]int, 0, len(members))
	for _, member := range members {
		id, err := strconv.Atoi(member)
		if err != nil {
			return nil, fmt.Errorf("parse indexed movie ID %q: %w", member, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
