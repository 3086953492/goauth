package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/database"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"github.com/3086953492/gokit/jwt"
	"gorm.io/driver/mysql"

	"goauth/initialize"
	"goauth/models"
)

func main() {
	// 初始化配置
	if err := config.InitConfig("./configs"); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	cfg := config.GetGlobalConfig()

	// 初始化日志
	if err := logger.InitWithConfig(cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	dsn := database.BuildMySQLDSN(cfg.Database)
	if err := database.InitDBWithDialector(mysql.Open(dsn)); err != nil {
		errors.Internal().Msg("初始化数据库失败").Err(err).Field("dsn", dsn).Log()
		return
	}

	models := []any{
		models.User{},
	}

	if err := database.AutoMigrateModels(models...); err != nil {
		errors.Internal().Msg("自动迁移数据库失败").Err(err).Field("models", models).Log()
		return
	}

	// 初始化 Redis
	if err := redis.InitRedisWithConfig(cfg.Redis); err != nil {
		errors.Internal().Msg("初始化 Redis 失败").Err(err).Field("cfg.Redis", cfg.Redis).Log()
		return
	}

	// 初始化缓存
	if err := cache.InitCache(); err != nil {
		errors.Internal().Msg("初始化缓存失败").Err(err).Log()
		return
	}

	if err := jwt.InitJWT(cfg.JWT); err != nil {
		errors.Internal().Msg("初始化 JWT 失败").Err(err).Log()
		return
	}

	container := initialize.NewContainer()

	initialize.InitValidator(container)

	// 获取端口号，优先使用命令行参数
	port := cfg.Server.Port
	if len(os.Args) > 1 {
		argPort, err := strconv.Atoi(os.Args[1])
		if err == nil {
			port = argPort
		}
	}

	// 初始化 Gin 路由，传入容器
	r := initialize.InitRouters(container)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		errors.Internal().Msg("启动服务失败").Err(err).Field("port", port).Log()
	}

	defer func() {
		if err := database.CloseDB(); err != nil {
			errors.Internal().Msg("关闭数据库失败").Err(err).Log()
		}
	}()
}
