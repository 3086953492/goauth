package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
	mgr, err := config.NewManager(
		config.WithConfigDir("./configs"),
		config.WithMode("debug"),
	)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg, err := mgr.Load(context.Background())
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

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
	redisMgr := redis.NewManager(
		redis.WithAddress(cfg.Redis.Host+":"+strconv.Itoa(cfg.Redis.Port)),
		redis.WithPassword(cfg.Redis.Password),
		redis.WithDB(cfg.Redis.DB),
	)

	if err := redisMgr.Connect(context.Background()); err != nil {
		errors.Internal().Msg("连接 Redis 失败").Err(err).Log()
		return
	}

	if !redisMgr.IsConnected() {
		errors.Internal().Msg("Redis 连接失败").Log()
		return
	}

	// 初始化缓存
	cacheMgr, err := cache.NewManager(redisMgr,
		cache.WithDefaultTTL(5*time.Minute),
		cache.WithLocalCache(true),
	)
	if err != nil {
		errors.Internal().Msg("初始化缓存失败").Err(err).Log()
		return
	}
	defer cacheMgr.Close()

	jwtMgr, err := jwt.NewManager(jwt.WithSecret(cfg.AuthToken.Secret),
		jwt.WithIssuer(cfg.AuthToken.Issuer),
		jwt.WithAccessTTL(cfg.AuthToken.AccessExpire),
		jwt.WithRefreshTTL(cfg.AuthToken.RefreshExpire),)
	if err != nil {
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

	container := initialize.NewContainer(dbManager.DB(), storageManager, validatorManager, redisMgr, cacheMgr, jwtMgr, &cfg)

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

	defer func() {
		if err := redisMgr.Close(); err != nil {
			errors.Internal().Msg("关闭 Redis 失败").Err(err).Log()
		}
	}()

	defer func() {
		if err := cacheMgr.Close(); err != nil {
			errors.Internal().Msg("关闭缓存失败").Err(err).Log()
		}
	}()
}
