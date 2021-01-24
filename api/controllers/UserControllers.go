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
		return "", "", err
		g
	}

	hubName, ok := hubForm["hubname"]
	username, ok2 := hubForm["username"]

	if !ok || !ok2 {
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

// JoinHub returns a cookie given a hub name and a player name
func (uc *UserController) JoinHub(w http.ResponseWriter, r *http.Request) {
	var claims security.PlayerClaims
	var ok bool
	hubName, playerName, errForm := getHubName(w, r)

	tokenString, errToken := extractToken(r)

	if errToken != nil {
		ok = sendMessageIfError(errForm, w, http.StatusUnprocessableEntity, "Player name can only have alphanumeric characters")
		if !ok {
			return
		}

		ok = sendMessageIfError(uc.checkInput(playerName), w, http.StatusUnprocessableEntity, "Player name can only have alphanumeric characters")
		if !ok {
			return
		}

		tokenString, errToken = security.GetToken(playerName, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))
		ok = sendMessageIfError(errToken, w, http.StatusInternalServerError, "")
		if !ok {
			return
		}

	} else {
		token, errToken := security.VerifyToken(tokenString, config.JWTSecret)

		if errToken != nil {
			tokenString, errToken = security.GetToken(playerName, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))

			ok = sendMessageIfError(errToken, w, http.StatusInternalServerError, "")
			if !ok {
				return
			}
		} else {
			claims, errToken = security.DecipherClaims(token)
			ok = sendMessageIfError(errToken, w, http.StatusInternalServerError, "")
			if !ok {
				return
			}

			if claims.HubName != hubName {
				if uc.checkInput(playerName) != nil {
					playerName = claims.PlayerName
				}

				tokenString, errToken = security.GetToken(playerName, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))
				ok = sendMessageIfError(errToken, w, http.StatusInternalServerError, "")
				if !ok {
					return
				}
			}
		}
	}

	errJSON := json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	ok = sendMessageIfError(errJSON, w, http.StatusInternalServerError, "")
	if !ok {
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

	ok = sendMessageIfError(err, w, http.StatusForbidden, "")
	if !ok {
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
	hubName, playerName, err := getHubName(w, r)

	if err != nil {
		uc.logger.LogChan <- err.Error()
		return
	}

	ok := sendMessageIfError(uc.checkInput(hubName), w, http.StatusUnprocessableEntity, "Hub name can only have alphanumeric characters")
	if !ok {
		return
	}

	ok = sendMessageIfError(uc.checkInput(playerName), w, http.StatusUnprocessableEntity, "Player name can only have alphanumeric characters")
	if !ok {
		return
	}

	err = uc.hubManager.RegisterHub(hubName)

	ok = sendMessageIfError(err, w, http.StatusInternalServerError, "")
	if !ok {
		return
	}

	tokenString, err := security.GetToken(playerName, hubName, config.JWTSecret, 0, int64(time.Duration(time.Minute*15)))

	ok = sendMessageIfError(err, w, http.StatusInternalServerError, "")
	if !ok {
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	ok = sendMessageIfError(err, w, http.StatusInternalServerError, "")
	if !ok {
		return
	}
}
