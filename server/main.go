package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/database"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/storage"
	"github.com/3086953492/gokit/storage/provider_aliyunoss"
	"github.com/3086953492/gokit/validator"
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
	dbManager, err := database.NewManager(mysql.Open(dsn))
	if err != nil {
		errors.Internal().Msg("初始化数据库失败").Err(err).Field("dsn", dsn).Log()
		return
	}
	defer dbManager.Close()

	models := []any{
		models.User{},
		models.OAuthClient{},
		models.OAuthAuthorizationCode{},
		models.OAuthAccessToken{},
		models.OAuthRefreshToken{},
	}

	if err := dbManager.AutoMigrate(models...); err != nil {
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

	if err := jwt.InitJWT(cfg.AuthToken); err != nil {
		errors.Internal().Msg("初始化 JWT 失败").Err(err).Log()
		return
	}

	store, err := provider_aliyunoss.New(provider_aliyunoss.Config(cfg.AliyunOSS))
	if err != nil {
		errors.Internal().Msg("初始化 OSS 失败").Err(err).Log()
		return
	}

	storageManager, err := storage.NewManager(storage.WithStore(store))
	if err != nil {
		errors.Internal().Msg("初始化文件管理器失败").Err(err).Log()
		return
	}

	response.Init(
		// 根据运行环境决定是否显示错误详情
		response.WithShowErrorDetail(cfg.Server.Mode == "debug"),
	)

	cookie.Init(cfg.Server.Mode != "debug")

	validatorManager, err := validator.New()
	if err != nil {
		errors.Internal().Msg("初始化验证器失败").Err(err).Log()
		return
	}

	container := initialize.NewContainer(dbManager.DB(), storageManager, validatorManager)

	if err := initialize.RegisterValidations(container); err != nil {
		errors.Internal().Msg("注册自定义验证规则失败").Err(err).Log()
		return
	}

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
		if err := dbManager.Close(); err != nil {
			errors.Internal().Msg("关闭数据库失败").Err(err).Log()
		}
	}()
}
