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
	Port        = "43389"
)

func main() {
	osc := vrcinput.NewOscClient(vrcinput.DefaultAddr, vrcinput.DefaultPort)

	t, err := template.ParseFiles("root.html", "view.html")
	if err != nil {
		log.Fatal(err)
	}

	h := live.NewHandler(live.WithTemplateRenderer(t))

	// Client side events.

	// Jump event.
	h.HandleEvent(Jump, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Jump()

		return 1, nil
	})

	// Move forward event.
	h.HandleEvent(MoveForward, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Move(vrcinput.MoveForward, true)
		time.Sleep(1 * time.Second)
		osc.Move(vrcinput.MoveForward, false)

		return 1, nil
	})

	// Move back event.
	h.HandleEvent(MoveBack, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Move(vrcinput.MoveBackward, true)
		time.Sleep(1 * time.Second)
		osc.Move(vrcinput.MoveBackward, false)

		return 1, nil
	})

	// Move left event.
	h.HandleEvent(MoveLeft, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Move(vrcinput.MoveLeft, true)
		time.Sleep(1 * time.Second)
		osc.Move(vrcinput.MoveLeft, false)

		return 1, nil
	})

	// Move right event.
	h.HandleEvent(MoveRight, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Move(vrcinput.MoveRight, true)
		time.Sleep(1 * time.Second)
		osc.Move(vrcinput.MoveRight, false)

		return 1, nil
	})

	// Look left event.
	h.HandleEvent(LookLeft, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Look(vrcinput.LookLeft, true)
		time.Sleep(250 * time.Millisecond)
		osc.Look(vrcinput.LookLeft, false)

		return 1, nil
	})

	// Look right event.
	h.HandleEvent(LookRight, func(ctx context.Context, s live.Socket, _ live.Params) (interface{}, error) {

		osc.Look(vrcinput.LookRight, true)
		time.Sleep(250 * time.Millisecond)
		osc.Look(vrcinput.LookRight, false)

		return 1, nil
	})

	// Run the server.
	go func() {
		http.Handle("/", live.NewHttpHandler(live.NewCookieStore("session-name", []byte("weak-secret")), h))
		http.Handle("/live.js", live.Javascript{})
		http.Handle("/auto.js.map", live.JavascriptMap{})
		err = http.ListenAndServe(":"+Port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server is running on http://localhost:" + Port)
	select {}

}
