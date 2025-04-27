package dependencies

import (
	"github.com/uptrace/bun"
	"taskapp/internal/cache"
)

type Dependencies struct {
	DB    *bun.DB
	Redis *cache.RedisCache
}
