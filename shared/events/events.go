package events

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/go-redis/redis/v8"
)

var ExternalSeedCooldown time.Duration = 1000 * time.Millisecond

type EventBus interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Subscribe(ctx context.Context, topic string, handler func(context.Context, []byte)) error
}

func NewRedisConnection(host, port, password string, dbNumber int) (*redis.Client, error) {
	// Створюємо клієнт Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       dbNumber,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

func NewDefaultRedisConnection() (*redis.Client, error) {
	dbNum, err := envutil.GetIntFromENV("REDIS_DB_NUMBER")
	if err != nil {
		return nil, fmt.Errorf("failed to get db number value: %s", err.Error())
	}

	return NewRedisConnection(
		envutil.GetStringFromENV("REDIS_HOST"),
		envutil.GetStringFromENV("REDIS_PORT"),
		envutil.GetStringFromENV("REDIS_PASSWORD"),
		dbNum,
	)
}

func NewEventBus(rdb *redis.Client) EventBus {
	return &redisPubSub{rdb: rdb}
}

type redisPubSub struct {
	rdb *redis.Client
}

func (r *redisPubSub) Publish(ctx context.Context, topic string, payload []byte) error {
	return r.rdb.Publish(ctx, topic, payload).Err()
}
func (r *redisPubSub) Subscribe(ctx context.Context, topic string, handler func(context.Context, []byte)) error {
	sub := r.rdb.Subscribe(ctx, topic)
	_, err := sub.Receive(ctx)
	if err != nil {
		return err
	}

	ch := sub.Channel()
	go func() {
		for msg := range ch {
			handler(contextutil.SetTraceID(context.Background()), []byte(msg.Payload))
		}
	}()
	return nil
}

func ToByteConvertor[T fmt.Stringer](entity T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(entity)

	if err != nil {
		return []byte{}, fmt.Errorf("failed to encode %s to bytes: %w", entity, err)
	}

	data := buf.Bytes()
	return data, nil
}
func FromByteConvertor[T fmt.Stringer](bEntity []byte) (T, error) {
	var entity T
	buf := bytes.NewBuffer(bEntity)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&entity)

	if err != nil {
		return entity, fmt.Errorf("failed to decode bytes: %w", err)
	}

	return entity, nil
}
