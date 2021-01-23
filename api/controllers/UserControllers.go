package controllers

import (
	"davidmultiplayersnake/api/config"
	"davidmultiplayersnake/api/models"
	"davidmultiplayersnake/api/multiplayer"
	"davidmultiplayersnake/api/security"
	"davidmultiplayersnake/utils"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// UserController controller for user routes
type UserController struct {
	upgrader   websocket.Upgrader
	hubManager *multiplayer.HubManager
	exp        *regexp.Regexp
	logger     *utils.Logger
}

// NewUserController returns a new *UserController
func NewUserController(hubManager *multiplayer.HubManager, logger *utils.Logger) *UserController {
	return &UserController{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hubManager: hubManager,
		exp:        regexp.MustCompile(`[^A-Za-z0-9]`),
		logger:     logger,
	}
}

func (uc *UserController) checkInput(input string) error {
	matches := uc.exp.FindAllStringIndex(input, -1)

	if input == "" || len(matches) > 0 {
		return errors.New("Invalid input")
	}

	return nil
}

func getHubName(w http.ResponseWriter, r *http.Request) (string, string, error) {
	var hubForm map[string]string

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&hubForm)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return "", "", err
	}

	hubName, ok := hubForm["hubname"]
	username, ok2 := hubForm["username"]

	if !ok || !ok2 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return "", "", errors.New("Hub name not found")
	}

	return hubName, username, nil
}

func extractToken(r *http.Request) (token string, err error) {
	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")

	if len(strArr) != 2 || strArr[0] != "Bearer" {
		err = errors.New("Invalid Token")
		return
	}

	token = strArr[1]

	return
}

// FetchHubs fetches all hubs and serializes them into JSON
func (uc *UserController) FetchHubs(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	w.Header().Add("Content-Type", "application/json")

	hubList := uc.hubManager.HubList()

	err := encoder.Encode(hubList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// JoinHub returns a cookie given a hubname and a username
func (uc *UserController) JoinHub(w http.ResponseWriter, r *http.Request) {
	var claims security.PlayerClaims

	hubName, username, err := getHubName(w, r)

	if err != nil {
		uc.logger.LogChan <- err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if uc.checkInput(username) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Player name can only have alphanumeric characters"))
		return
	}

	tokenString, err := extractToken(r)

	if err != nil {
		tokenString, err = security.GetToken(username, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))
	} else {
		token, err := security.VerifyToken(tokenString, config.JWTSecret)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			claims, err = security.DecipherClaims(token)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if claims.HubName != hubName {
			tokenString, err = security.GetToken(username, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))
		}
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	if err != nil {
		uc.logger.LogChan <- err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// StartGame starts the game
func (uc *UserController) StartGame(w http.ResponseWriter, r *http.Request) {
	playerClaims, ok := r.Context().Value(security.PlayerClaimsType).(security.PlayerClaims)

	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	h := uc.hubManager.GetHub(playerClaims.HubName)

	if h == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	connection, err := uc.upgrader.Upgrade(w, r, nil)

	if err != nil {
		uc.logger.LogChan <- err.Error()
		w.WriteHeader(http.StatusForbidden)
		return
	}

	client := multiplayer.NewClient(connection, h, nil)

	client.Player = models.NewPlayer(playerClaims.PlayerName, models.Position{X: rand.Intn(1999) + 1, Y: rand.Intn(999) + 1}, models.Direction(rand.Intn(4)+1))
	client.Player.ID = playerClaims.PlayerID
	client.Player.Score = playerClaims.PlayerScore

	h.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

// CreateHub is the controller to create a new Hub
func (uc *UserController) CreateHub(w http.ResponseWriter, r *http.Request) {
	hubName, username, err := getHubName(w, r)

	if err != nil {
		uc.logger.LogChan <- err.Error()
		return
	}

	if uc.checkInput(hubName) != nil || uc.checkInput(hubName) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Hub name and player name can only have alphanumeric characters"))
		return
	}

	err = uc.hubManager.RegisterHub(hubName)

	if err != nil {
		uc.logger.LogChan <- err.Error()

		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != nil {
		uc.logger.LogChan <- err.Error()

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenString, err := security.GetToken(username, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))

	if err != nil {
		uc.logger.LogChan <- err.Error()

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	if err != nil {
		uc.logger.LogChan <- err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
