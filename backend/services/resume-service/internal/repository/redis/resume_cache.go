package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"resume-service/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type ResumeCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewResumeCache(client *redis.Client) *ResumeCache {
	return &ResumeCache{
		client: client,
		ttl:    15 * time.Minute,
	}
}

func (c *ResumeCache) Get(ctx context.Context, key string) (*entity.Resume, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var resume entity.Resume
	if err := json.Unmarshal([]byte(data), &resume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resume: %w", err)
	}

	return &resume, nil
}

func (c *ResumeCache) Set(ctx context.Context, key string, resume *entity.Resume) error {
	data, err := json.Marshal(resume)
	if err != nil {
		return fmt.Errorf("failed to marshal resume: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *ResumeCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (c *ResumeCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	return iter.Err()
}
