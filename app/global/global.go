package global

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"zhihu/app/internal/model/config"
)

var (
	Config *config.Config
	Logger *zap.Logger
	Mysql  *sql.DB
	Redis  *redis.Client
)
