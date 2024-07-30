package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	vrcinput "github.com/Jilwer/VRChatOscInput"
	"github.com/jfyne/live"
)

const (
	MoveForward = "up"
	MoveBack    = "down"
	MoveLeft    = "left"
	MoveRight   = "right"
	Jump        = "jump"
	LookLeft    = "look-left"
	LookRight   = "look-right"
	Port        = "8080"
	Domain      = "remote.ubel.org"
	Send        = "send"
)

func main() {
	osc := initializeOscClient()
	t := initializeTemplates()
	h := initializeHandlers(&osc, t)

	go runServer(h)

	fmt.Println("Server is running on http://localhost:" + Port)
	select {}
}

func runServer(h live.Handler) {
	http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Initializers

func initializeOscClient() vrcinput.Client {
	osc := vrcinput.NewOscClient(vrcinput.DefaultAddr, vrcinput.DefaultPort)
	go func() {
		for {
			osc.Chat("Control me at: "+Domain, true, false)
			time.Sleep(10 * time.Second)
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
		return handleSendEvent(osc, p)
	})

	h.HandleEvent(Jump, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {
		return handleJumpEvent(osc)
	})

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

	return h
}

// Event handlers

func handleSendEvent(osc *vrcinput.Client, p live.Params) (interface{}, error) {
	msg := p.String("message")
	if msg == "" {
		return 1, nil
	}

	if len(msg) > 143 {
		msg = msg[:143]
	}

	osc.Chat(msg, true, false)
	return 1, nil
}

func handleJumpEvent(osc *vrcinput.Client) (interface{}, error) {
	osc.Jump()
	return 1, nil
}

func handleMoveEvent(osc *vrcinput.Client, direction vrcinput.MoveDirection) (interface{}, error) {
	osc.Move(direction, true)
	time.Sleep(1 * time.Second)
	osc.Move(direction, false)
	return 1, nil
}

func handleLookEvent(osc *vrcinput.Client, direction vrcinput.LookDirection) (interface{}, error) {
	osc.Look(direction, true)
	time.Sleep(250 * time.Millisecond)
	osc.Look(direction, false)
	return 1, nil
}
