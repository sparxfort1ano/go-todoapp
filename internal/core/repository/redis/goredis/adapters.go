package redispool

import (
	"errors"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/redis"
)

type goredisStringCmd struct {
	*goredis.StringCmd
}

func (c goredisStringCmd) Bytes() ([]byte, error) {
	result, err := c.StringCmd.Bytes()
	if err != nil {
		return nil, mapError(err)
	}

	return result, nil
}

type goredisStatusCmd struct {
	*goredis.StatusCmd
}

type goredisIntCmd struct {
	*goredis.IntCmd
}

// mapError converts errors specific to this Redis client into domain-specific errors.
func mapError(err error) error {
	switch {
	case errors.Is(err, goredis.Nil):
		return redis.ErrNotFound
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		redis.ErrUnknown,
	)
}
