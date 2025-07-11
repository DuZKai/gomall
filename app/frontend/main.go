// Code generated by hertz generator.

package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/logger/accesslog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/sessions"
	"github.com/hertz-contrib/sessions/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
	"gomall/app/frontend/biz/router"
	"gomall/app/frontend/conf"
	"gomall/app/frontend/infra/rpc"
	"gomall/app/frontend/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()
	// init dal
	// dal.Init()
	rpc.InitClient()
	address := conf.GetConf().Hertz.Address
	h := server.New(server.WithHostPorts(address))

	registerMiddleware(h)

	// add a ping route to test
	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	router.GeneratedRegister(h)
	// 加载模板文件
	h.LoadHTMLGlob("template/*")

	//h.Static("/static", "./")
	h.StaticFS("/static", &app.FS{
		Root:               "./",
		GenerateIndexPages: true,
		PathRewrite:        nil,
	})

	h.GET("/about", func(c context.Context, ctx *app.RequestContext) {
		ctx.HTML(consts.StatusOK, "about", utils.H{"title": "About"})
	})

	h.GET("/sign-in", func(c context.Context, ctx *app.RequestContext) {
		data := utils.H{
			"title": "Sign In",
			"next":  ctx.Query("next"),
		}
		ctx.HTML(consts.StatusOK, "sign-in", data)
	})

	h.GET("/sign-up", func(c context.Context, ctx *app.RequestContext) {
		ctx.HTML(consts.StatusOK, "sign-up", utils.H{"title": "Sign Up"})
	})

	h.Spin()
}

func registerMiddleware(h *server.Hertz) {
	store, err := redis.NewStore(10, "tcp", conf.GetConf().Redis.Address, conf.GetConf().Redis.Password, []byte(os.Getenv("SESSION_SECRET")))
	
	if err != nil {
		panic(err)
	}
	h.Use(sessions.New("gomall-session", store))

	// log
	logger := hertzlogrus.NewLogger()
	hlog.SetLogger(logger)
	hlog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Hertz.LogFileName,
			MaxSize:    conf.GetConf().Hertz.LogMaxSize,
			MaxBackups: conf.GetConf().Hertz.LogMaxBackups,
			MaxAge:     conf.GetConf().Hertz.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	hlog.SetOutput(asyncWriter)
	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		asyncWriter.Sync()
	})

	// pprof
	if conf.GetConf().Hertz.EnablePprof {
		pprof.Register(h)
	}

	// gzip
	if conf.GetConf().Hertz.EnableGzip {
		h.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	// access log
	if conf.GetConf().Hertz.EnableAccessLog {
		h.Use(accesslog.New())
	}

	// recovery
	h.Use(recovery.Recovery())

	// cores
	h.Use(cors.Default())

	middleware.Register(h)
}
