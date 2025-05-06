package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisMemory struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisMemory(addr, pass string, db int) *RedisMemory {
	return &RedisMemory{client: redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db}), ctx: context.Background()}
}

func (r *RedisMemory) key(userID int64) string {
	return fmt.Sprintf("chat_history:%d", userID)
}

func (r *RedisMemory) Append(userID int64, msg ChatMessage) error {
	data, _ := json.Marshal(msg)
	return r.client.RPush(r.ctx, r.key(userID), data).Err()
}

func (r *RedisMemory) Get(userID int64) ([]ChatMessage, error) {
	vals, err := r.client.LRange(r.ctx, r.key(userID), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	out := make([]ChatMessage, len(vals))
	for i, v := range vals {
		_ = json.Unmarshal([]byte(v), &out[i])
	}
	return out, nil
}

func (r *RedisMemory) Truncate(userID int64, limit int) error {
	// оставить последние limit элементов
	if limit <= 0 {
		return r.client.Del(r.ctx, r.key(userID)).Err()
	}
	return r.client.LTrim(r.ctx, r.key(userID), int64(-limit), -1).Err()
}
