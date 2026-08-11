// Package cached consists of a set of adapters that adapt cache to the outcoming port (ports/out/repository),
// interacting with models responsible for the actual network interaction with the caching database (i.e Redis).
package cached

import (
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/redis"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// CachedRepository is a caching adapter that wraps the underlying repository (e.g. postgres).
//
// Caching strategies:
//   - Single task (key "task:<id>"): cache-aside using GET/SET.
//   - List of tasks (key "tasks:<userID>" / "tasks:all"): cache-aside via a hash,
//     where the hash’s fields encode pagination (<limit>:<offset>).
//
// Invalidation occurs on any mutation (i.e. save, update, delete).
//
// Cache errors are logged and do not interrupt the request — this ensures graceful degradation to the primary storage.
type CachedRepository struct {
	pool           redis.Pool
	mainRepository repository.TasksRepository
}

func NewCachedRepository(
	pool redis.Pool,
	mainRepository repository.TasksRepository,
) *CachedRepository {
	return &CachedRepository{
		pool:           pool,
		mainRepository: mainRepository,
	}
}
