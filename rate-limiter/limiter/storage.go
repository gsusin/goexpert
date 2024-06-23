package limiter

import (
	"context"
	redis "github.com/redis/go-redis/v9"
)

type MemoryStorage struct {
	access map[string]int
}

func NewMemoryStorage() MemoryStorage {
	return MemoryStorage{
		make(map[string]int),
	}
}

func (s *MemoryStorage) SetZero(ip string) {
	s.access[ip] = 0
}

func (s *MemoryStorage) Increment(ip string) {
	s.access[ip]++
}

func (s *MemoryStorage) Get(ip string) int {
	return s.access[ip]
}

type RedisStorage struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisStorage() RedisStorage {
	return RedisStorage{
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0, // use default DB
		}),
		context.Background(),
	}
}

func (s *RedisStorage) SetZero(ip string) {
	err := s.rdb.Set(s.ctx, "access", 0, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (s *RedisStorage) Increment(ip string) {
	val := s.Get(ip)
	err := s.rdb.Set(s.ctx, "access", val+1, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (s *RedisStorage) Get(ip string) int {
	val, err := s.rdb.Get(s.ctx, "access").Int()
	if err != nil {
		panic(err)
	}
	return val
}
