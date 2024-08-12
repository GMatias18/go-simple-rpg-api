package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type PlayerRequest struct {
	Nickname string `json:"nickname"`
	Life     int    `json:"life"`
	Attack   int    `json:"attack"`
}

type PlayerResponse struct {
	Message string `json:"message"`
}

type Enemy struct {
	Nickname string `json:"nickname"`
	Life     int    `json:"life"`
	Attack   int    `json:"attack"`
}

type Battle struct {
	ID         string `json:"id"`
	Enemy      string `json:"enemy"`
	Player     string `json:"player"`
	DiceThrown int    `json:"dice_thrown"`
}

var players []PlayerRequest
var enemies []Enemy
var battles []Battle

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/player", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddPlayer(w, r)
		case http.MethodGet:
			LoadPlayers(w, r)
		}
	})
	mux.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			DeletePlayer(w, r)
		case http.MethodGet:
			LoadPlayerByNickname(w, r)
		case http.MethodPut:
			SavePlayer(w, r)
		}
	})

	mux.HandleFunc("/enemy", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddEnemy(w, r)
		case http.MethodGet:
			LoadEnemies(w, r)
		}
	})
	mux.HandleFunc("/enemy/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			DeleteEnemy(w, r)
		case http.MethodGet:
			LoadEnemyByNickname(w, r)
		case http.MethodPut:
			UpdateEnemy(w, r)
		}
	})

	mux.HandleFunc("/battle", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateBattle(w, r)
		case http.MethodGet:
			LoadBattles(w, r)
		}
	})

	fmt.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", mux)
}

// Implementações de AddPlayer, LoadPlayers, DeletePlayer, etc. não são repetidas aqui para economia de espaço.

func LoadBattles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battles)
}
