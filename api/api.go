package api

import (
	"bufio"
	"davidmultiplayersnake/api/controllers"
	"davidmultiplayersnake/api/middlewares"
	"davidmultiplayersnake/api/multiplayer"
	"davidmultiplayersnake/utils"
	"net/http"
	"os"
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
	mux.HandleFunc("/game", middlewares.MethodMiddleware(middlewares.JWTMiddleware(uc.StartGame, true), http.MethodGet))

	// Frontend routes
	fsGame := http.StripPrefix("/play/", http.FileServer(http.Dir("../gamefiles")))
	fsHome := http.FileServer(http.Dir("../home"))

	mux.HandleFunc("/play/", middlewares.MethodMiddleware(middlewares.JWTMiddleware(fsGame.ServeHTTP, true), http.MethodGet))

	mux.HandleFunc("/", middlewares.MethodMiddleware(middlewares.JWTMiddleware(fsHome.ServeHTTP, false), http.MethodGet))

	// HTTP Server
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Duration(3 * time.Second),
		WriteTimeout: time.Duration(2 * time.Second),
		Handler:      mux,
	}

	// Let hubs be created concurrently
	go hm.Run()
	// Let logs appear concurrently
	go l.Logs()

	s.ListenAndServe()
}
