package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/websocket/v2"
	"github.com/gonevo/afterall"
	_ "github.com/gonevo/automaxprocs"
	"github.com/gonevo/logium/internal/config"
	"github.com/gonevo/logium/internal/logline"
	"github.com/gonevo/logium/internal/source"
	"github.com/gonevo/logium/internal/ticker"
	"github.com/gonevo/logium/internal/welcome"
	"log"
	"os"
	"syscall"
)

//go:embed static/index.html
var template string

var sources = source.New()
var welcomeLines *welcome.Welcome

// websocket clients
type client struct{}

var clients = make(map[*websocket.Conn]client)

// channels
var register = make(chan *websocket.Conn, 100)
var unregister = make(chan *websocket.Conn, 100)
var broadcast = make(chan *logline.LogLine, 100)

func main() {

	var i bool
	flag.BoolVar(&i, "init", false, "Create new config.json")
	flag.Parse()
	if i {
		config.Template()
		os.Exit(0)
	}

	appConfig := config.Get()

	if len(appConfig.Sources) == 0 {
		log.Printf("\nSources are not configured. Please fill %s", config.FILE)
		os.Exit(1)
	}

	welcomeLines = welcome.New(appConfig.Server.Last)

	initSources(appConfig.Sources, welcomeLines)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	if len(appConfig.Server.Auth) != 0 {
		app.Use(basicauth.New(basicauth.Config{
			Users: appConfig.Server.Auth,
		}))
	}

	app.Use(favicon.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			_ = c.Close()
		}()

		register <- c

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return
			}
		}
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return ctx.SendString(template)
	})

	go func() {
		_ = app.Listen(fmt.Sprintf(":%d", appConfig.Server.Port))
	}()

	go runBroadcast()

	if tick := os.Getenv("TICKER"); tick != "" {
		go func() {
			ticker.Tick(appConfig)
		}()
	}

	for path, s := range sources.Get() {
		go initSourceListener(path, s)
	}

	afterall.
		I().
		On(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP).
		HaveToCall(func() {
			sources.Stop()
			os.Exit(0)
		}).
		Wait()
}

func initSources(inputs []config.Source, w *welcome.Welcome) {
	for _, input := range inputs {
		sources.Init(input.Path, input.Name, w)
	}
}

func initSourceListener(path string, s *source.Source) {
	for l := range s.Tail.Lines {
		logLine := logline.New(s.Name, l.Text)
		welcomeLines.Append(path, logLine)
		broadcast <- logLine
	}
}

func sendMessage(connection *websocket.Conn, logLine *logline.LogLine) {
	if err := connection.WriteMessage(websocket.TextMessage, logLine.ToBytes()); err != nil {
		unregister <- connection
		_ = connection.WriteMessage(websocket.CloseMessage, []byte{})
		_ = connection.Close()
	}
}

func runBroadcast() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			for _, lines := range welcomeLines.Get() {
				for _, logLine := range lines {
					sendMessage(connection, logLine)
				}
			}

		case logLine := <-broadcast:
			for connection := range clients {
				sendMessage(connection, logLine)
			}

		case connection := <-unregister:
			delete(clients, connection)
		}
	}
}
