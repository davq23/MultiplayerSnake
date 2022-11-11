package api

import (
	"bufio"
	"context"
	"davidmultiplayersnake/api/controllers"
	"davidmultiplayersnake/api/middlewares"
	"davidmultiplayersnake/api/multiplayer"
	"davidmultiplayersnake/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run starts backend api
func Run() {
	mux := http.NewServeMux()

	l := utils.NewLogger(bufio.NewWriterSize(os.Stdout, 256))
	hm := multiplayer.NewHubManager(l)
	uc := controllers.NewUserController(hm, l)

	// Backend routes
	mux.HandleFunc("/hubs/create", middlewares.MethodMiddleware(middlewares.JWTMiddleware(uc.CreateHub, false), http.MethodPost))
	mux.HandleFunc("/hubs/join", middlewares.MethodMiddleware(middlewares.JWTMiddleware(uc.JoinHub, false), http.MethodPost))
	mux.HandleFunc("/hubs", middlewares.MethodMiddleware(middlewares.JWTMiddleware(uc.FetchHubs, false), http.MethodGet))

	// Websocket route
	mux.HandleFunc("/game", middlewares.MethodMiddleware(middlewares.WSJWTMiddleware(uc.StartGame), http.MethodGet))

	// Frontend routes
	fsGame := http.StripPrefix("/play/", http.FileServer(http.Dir("./gamefiles")))
	fsHome := http.FileServer(http.Dir("./home"))

	mux.HandleFunc("/play/", middlewares.MethodMiddleware(fsGame.ServeHTTP, http.MethodGet))

	mux.HandleFunc("/", middlewares.MethodMiddleware(fsHome.ServeHTTP, http.MethodGet))

	// HTTP Server
	s := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		ReadTimeout:  time.Duration(3 * time.Second),
		WriteTimeout: time.Duration(2 * time.Second),
		Handler:      mux,
	}

	// Let hubs be created concurrently
	go hm.Run()
	// Let logs appear concurrently
	go l.Logs()

	go s.ListenAndServe()

	// Graceful termination
	sigChan := make(chan os.Signal, 5)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan

	log.Println("Received terminate, graceful shutdown", sig)

	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancelFunc()

	defer func() {
		// Stop
		l.Running = false
		hm.Running = false
		s.Shutdown(tc)
	}()
}
