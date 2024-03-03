package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/funkymcb/guillocut/components"
	"github.com/funkymcb/guillocut/config"
	"github.com/funkymcb/guillocut/gintemplrenderer"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	engine := gin.New()

	logger, _ := zap.NewProduction()
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	if err := engine.SetTrustedProxies(nil); err != nil {
		log.Fatalln("could not set gin trusted proxies option", err)
	}

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", components.Home())
	})

	cfg, err := config.Get()
	if err != nil {
		slog.Error("error processing config", "message", err.Error())
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := engine.Run(addr); err != nil {
		log.Fatalln("error executing gin server", err)
	}
}
