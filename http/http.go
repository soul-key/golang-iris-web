package http

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
	"context"
	"strconv"
	"sync"
	"time"
)

func RunIris(port int) {
	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	// 优雅的关闭程序
	serverWG := new(sync.WaitGroup)
	defer serverWG.Wait()

	iris.RegisterOnInterrupt(func() {
		serverWG.Add(1)
		defer serverWG.Done()

		timeout := time.Second * 5
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// 关闭所有主机
		app.Shutdown(ctx)
	})

	// 注册路由
	innerRoute(app)

	// server配置
	c := iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           true,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
		RemoteAddrHeaders: 				   map[string]bool{"X-Real-Ip": true,"X-Forwarded-For": true},
	})

	app.Run(iris.Addr(":" + strconv.Itoa(port)), c)
}