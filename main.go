package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	vrcinput "github.com/Jilwer/VRChatOscInput"
	"github.com/jfyne/live"
)

const (
	Jump            = "jump"
	MoveForward     = "up"
	MoveForwardStop = "up-stop"
	MoveBack        = "down"
	MoveBackStop    = "down-stop"
	MoveLeft        = "left"
	MoveLeftStop    = "left-stop"
	MoveRight       = "right"
	MoveRightStop   = "right-stop"
	LookLeft        = "look-left"
	LookLeftStop    = "look-left-stop"
	LookRight       = "look-right"
	LookRightStop   = "look-right-stop"
	Send            = "send"
)

var (
	lastMessageTime   time.Time
	rateLimit         = 1 * time.Second
	mu                sync.Mutex
	sendStaticMessage = true
	UserConfig        *Config
)

func init() {
	var err error
	UserConfig, err = LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	osc := initializeOscClient()
	t := initializeTemplates()
	h := initializeHandlers(&osc, t)

	go runServer(h)

	fmt.Println("Server is running on http://localhost:" + UserConfig.Port)
	select {}
}

func runServer(h live.Handler) {
	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	err := http.ListenAndServe(":"+UserConfig.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Initializers

// create a global sendStaticMessage timer

var sendStaticMessageTimer = time.NewTicker(10 * time.Second)

func initializeOscClient() vrcinput.Client {
	osc := vrcinput.NewOscClient(vrcinput.DefaultAddr, vrcinput.DefaultPort)
	go func() {
		for range sendStaticMessageTimer.C {
			if sendStaticMessage {
				osc.Chat(UserConfig.StaticMessage, true, false)
			}
		}
	}()
	return osc
}
func initializeTemplates() *template.Template {
	t, err := template.ParseFiles("root.html", "view.html")
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func initializeHandlers(osc *vrcinput.Client, t *template.Template) live.Handler {
	h := live.NewHandler(live.WithTemplateRenderer(t))

	h.HandleEvent(Send, func(ctx context.Context, s live.Socket, p live.Params) (interface{}, error) {
		return handleChatEvent(osc, p)
	})

	h.HandleEvent(Jump, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleJumpEvent(osc)
	})

	// Input Start Events
	h.HandleEvent(MoveForward, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleMoveEvent(osc, vrcinput.MoveForward)
	})

	h.HandleEvent(MoveBack, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleMoveEvent(osc, vrcinput.MoveBackward)
	})

	h.HandleEvent(MoveLeft, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleMoveEvent(osc, vrcinput.MoveLeft)
	})

	h.HandleEvent(MoveRight, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleMoveEvent(osc, vrcinput.MoveRight)
	})

	h.HandleEvent(LookLeft, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleLookEvent(osc, vrcinput.LookLeft)
	})

	h.HandleEvent(LookRight, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleLookEvent(osc, vrcinput.LookRight)
	})

	// Input Stop Events
	h.HandleEvent(MoveForwardStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopMoveEvent(osc, vrcinput.MoveForward)
	})

	h.HandleEvent(MoveBackStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopMoveEvent(osc, vrcinput.MoveBackward)
	})

	h.HandleEvent(MoveLeftStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopMoveEvent(osc, vrcinput.MoveLeft)
	})

	h.HandleEvent(MoveRightStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopMoveEvent(osc, vrcinput.MoveRight)
	})

	h.HandleEvent(LookLeftStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopLookEvent(osc, vrcinput.LookLeft)
	})

	h.HandleEvent(LookRightStop, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleStopLookEvent(osc, vrcinput.LookRight)
	})

	return h
}

// Event handlers

func handleChatEvent(osc *vrcinput.Client, p live.Params) (interface{}, error) {

	if !UserConfig.ChatEnabled {
		return 1, nil
	}

	msg := p.String("message")
	if msg == "" {
		return 1, nil
	}

	if len(msg) > 143 {
		msg = msg[:143]
	}

	mu.Lock()
	defer mu.Unlock()

	if time.Since(lastMessageTime) < rateLimit {
		log.Println("Chat rate limit exceeded")
		return nil, fmt.Errorf("rate limit exceeded")
	}

	osc.Chat(msg, true, false)

	sendStaticMessageTimer.Reset(10 * time.Second)

	lastMessageTime = time.Now()
	return 1, nil
}

func handleJumpEvent(osc *vrcinput.Client) (interface{}, error) {
	osc.Jump()
	return 1, nil
}

func handleMoveEvent(osc *vrcinput.Client, direction vrcinput.MoveDirection) (interface{}, error) {
	osc.Move(direction, true)
	// log.Println("Received move event: ", direction)
	return 1, nil
}

func handleStopMoveEvent(osc *vrcinput.Client, direction vrcinput.MoveDirection) (interface{}, error) {
	osc.Move(direction, false)
	// log.Println("Received stop event: ", direction)
	return 1, nil
}

func handleLookEvent(osc *vrcinput.Client, direction vrcinput.LookDirection) (interface{}, error) {
	osc.Look(direction, true)
	return 1, nil
}

func handleStopLookEvent(osc *vrcinput.Client, direction vrcinput.LookDirection) (interface{}, error) {
	//log.Println("Received stop look event: ", direction)
	osc.Look(direction, false)
	return 1, nil
}
