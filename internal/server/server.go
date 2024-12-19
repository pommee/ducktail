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
	LogsCh    chan string
	IndexHTML string
	StyleCSS  string
	ScriptJS  string
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
	server.serveWebsite(opts)

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

func (server *API) serveWebsite(opts ServerOpts) {
	staticFiles := map[string]struct {
		content  string
		mimeType string
	}{
		"/":                     {opts.IndexHTML, "text/html"},
		"/static/css/style.css": {opts.StyleCSS, "text/css"},
		"/static/js/script.js":  {opts.ScriptJS, "application/javascript"},
	}

	for route, file := range staticFiles {
		server.router.GET(route, func(c *gin.Context) {
			c.Data(200, file.mimeType, []byte(file.content))
		})
	}
}
