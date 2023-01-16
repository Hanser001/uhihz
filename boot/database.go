package boot

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
	g "zhihu/app/global"
)

func MysqlSet() {
	config := g.Config.Database.Mysql

	db, err := sql.Open("mysql", config.GetDsn())
	if err != nil {
		g.Logger.Fatal("initialize mysql failed.", zap.Error(err))
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.GetConnMaxLifeTime())
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxIdleTime(config.GetConnMaxIdleTime())

	err = db.Ping()
	if err != nil {
		g.Logger.Fatal("connect to mysql failed", zap.Error(err))
	}

	g.Mysql = db

	g.Logger.Info("initialize mysql successfully!")
}

func RedisSet() {
	//config := g.Config.Database.Redis 用配置文件会报错空指针，先直连进行测试

	rdb := redis.NewClient(&redis.Options{
		Addr:     "39.101.68.42:6379",
		Password: "Wjj20040311",
		DB:       0,
		PoolSize: 10,
		Network:  "tcp",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		g.Logger.Fatal("connect to redis instance failed.", zap.Error(err))
	}

	g.Redis = rdb

	g.Logger.Info("initialize redis client successfully!")
}
