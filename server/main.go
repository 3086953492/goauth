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
	)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg, err := mgr.Load(context.Background())
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logMgr, err := logger.NewManager(logger.WithLevelString(cfg.Log.Level), logger.WithConsole(cfg.Server.Mode != "release"), logger.WithFile(logger.FileConfig{
		Filename:       cfg.Log.Filename,
		MaxSize:        cfg.Log.MaxSize,
		MaxAge:         cfg.Log.MaxAge,
		MaxBackups:     cfg.Log.MaxBackups,
		Compress:       cfg.Log.Compress,
		RotateStrategy: logger.RotateByDate,
	}))
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logMgr.Close()

	// 初始化数据库
	dsn := database.BuildMySQLDSN(cfg.Database)
	dbManager, err := database.NewManager(mysql.Open(dsn))
	if err != nil {
		logMgr.Error("初始化数据库失败", "dsn", dsn, "error", err)
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
		logMgr.Error("自动迁移数据库失败", "models", models, "error", err)
		return
	}

	// 初始化 Redis
	redisMgr := redis.NewManager(
		redis.WithAddress(cfg.Redis.Host+":"+strconv.Itoa(cfg.Redis.Port)),
		redis.WithPassword(cfg.Redis.Password),
		redis.WithDB(cfg.Redis.DB),
	)

	if err := redisMgr.Connect(context.Background()); err != nil {
		logMgr.Error("连接 Redis 失败", "error", err)
		return
	}

	if !redisMgr.IsConnected() {
		logMgr.Error("Redis 连接失败")
		return
	}

	// 初始化缓存
	cacheMgr, err := cache.NewManager(redisMgr,
		cache.WithDefaultTTL(5*time.Minute),
		cache.WithLocalCache(true),
	)
	if err != nil {
		logMgr.Error("初始化缓存失败", "error", err)
		return
	}
	defer cacheMgr.Close()

	jwtMgr, err := jwt.NewManager(jwt.WithSecret(cfg.AuthToken.Secret),
		jwt.WithIssuer(cfg.AuthToken.Issuer),
		jwt.WithAccessTTL(cfg.AuthToken.AccessExpire),
		jwt.WithRefreshTTL(cfg.AuthToken.RefreshExpire))
	if err != nil {
		logMgr.Error("初始化 JWT 失败", "error", err)
		return
	}

	store, err := provider_aliyunoss.New(provider_aliyunoss.Config(cfg.AliyunOSS))
	if err != nil {
		logMgr.Error("初始化 OSS 失败", "error", err)
		return
	}

	storageManager, err := storage.NewManager(storage.WithStore(store))
	if err != nil {
		logMgr.Error("初始化文件管理器失败", "error", err)
		return
	}

	response.Init(
		// 根据运行环境决定是否显示错误详情
		response.WithShowErrorDetail(cfg.Server.Mode == "debug"),
	)

	cookie.Init(cfg.Server.Mode != "debug")

	validatorManager, err := validator.New()
	if err != nil {
		logMgr.Error("初始化验证器失败", "error", err)
		return
	}

	container := initialize.NewContainer(dbManager.DB(), storageManager, validatorManager, redisMgr, cacheMgr, jwtMgr, logMgr, &cfg)

	if err := initialize.RegisterValidations(container); err != nil {
		logMgr.Error("注册自定义验证规则失败", "error", err)
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
		logMgr.Error("启动服务失败", "port", port, "error", err)
	}

	defer func() {
		if err := dbManager.Close(); err != nil {
			logMgr.Error("关闭数据库失败", "error", err)
		}
	}()

	defer func() {
		if err := redisMgr.Close(); err != nil {
			logMgr.Error("关闭 Redis 失败", "error", err)
		}
	}()

	defer func() {
		if err := cacheMgr.Close(); err != nil {
			logMgr.Error("关闭缓存失败", "error", err)
		}
	}()
}
