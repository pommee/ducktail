package server

import (
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	logs []string
)

type API struct {
	router *gin.Engine
	ws     *websocket.Conn
	logsCh <-chan string
}

type ServerOpts struct {
	LogsCh chan string
}

func (server *API) readLogs() {
	for logLine := range server.logsCh {
		logs = append(logs, logLine)
	}
}

func (server *API) create(opts ServerOpts) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	server.router = router
	server.logsCh = opts.LogsCh
}

func (server *API) Start(opts ServerOpts) {
	log.SetFlags(log.Flags() &^ (log.Ldate))

	server.create(opts)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.readLogs()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.router.Run(":8412")
	}()

	server.serve()
	server.serveWebsite()

	slog.Info("Started server at http://localhost:8412")
	wg.Wait()
}

func (server *API) serve() {
	server.router.GET("/logs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"logs": logs,
		})
	})

}

func (server *API) serveWebsite() {
	server.router.LoadHTMLFiles("website/index.html")
	server.router.Static("/static", "./website/static")

	server.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
