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
	Defense  int    `json:"defense"`
}

type PlayerResponse struct {
	Message string `json:"message"`
}

type Enemy struct {
	Nickname string `json:"nickname"`
	Life     int    `json:"life"`
	Attack   int    `json:"attack"`
	Defense  int    `json:"defense"`
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

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}
	for _, player := range players {
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname already exists"})
			return
		}
	}
	players = append(players, playerRequest)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playerRequest)
}

func LoadPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	for i, player := range players {
		if player.Nickname == nickname {
			players = append(players[:i], players[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname not found"})
}

func LoadPlayerByNickname(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	for _, player := range players {
		if player.Nickname == nickname {
			json.NewEncoder(w).Encode(player)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname not found"})
}

func SavePlayer(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}
	for i, player := range players {
		if player.Nickname == nickname {
			players[i] = playerRequest
			json.NewEncoder(w).Encode(playerRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Player not found"})
}

func AddEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var enemyRequest Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}
	for _, enemy := range enemies {
		if enemy.Nickname == enemyRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname already exists"})
			return
		}
	}
	rand.Seed(time.Now().UnixNano())
	enemyRequest.Life = rand.Intn(10) + 1
	enemyRequest.Attack = rand.Intn(10) + 1
	enemies = append(enemies, enemyRequest)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemyRequest)
}

func LoadEnemies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enemies)
}

func LoadEnemyByNickname(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	for _, enemy := range enemies {
		if enemy.Nickname == nickname {
			json.NewEncoder(w).Encode(enemy)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func UpdateEnemy(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	var enemyRequest Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}
	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			enemyRequest.Life = enemy.Life
			enemyRequest.Attack = enemy.Attack
			enemies[i] = enemyRequest
			json.NewEncoder(w).Encode(enemyRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy not found"})
}

func DeleteEnemy(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			enemies = append(enemies[:i], enemies[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func CreateBattle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var battleRequest struct {
		Enemy  string `json:"enemy"`
		Player string `json:"player"`
	}
	if err := json.NewDecoder(r.Body).Decode(&battleRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}
	var player *PlayerRequest
	var enemy *Enemy
	playerFound, enemyFound := false, false
	for i := range players {
		if players[i].Nickname == battleRequest.Player {
			player = &players[i]
			playerFound = true
			break
		}
	}
	for i := range enemies {
		if enemies[i].Nickname == battleRequest.Enemy {
			enemy = &enemies[i]
			enemyFound = true
			break
		}
	}
	if !playerFound || !enemyFound {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player or Enemy not found"})
		return
	}
	if player.Life <= 0 || enemy.Life <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "One of the combatants is dead, battle cannot proceed"})
		return
	}
	diceThrown := rand.Intn(6) + 1
	playerDamage := max(0, player.Attack - enemy.Defense + diceThrown)
	enemyDamage := max(0, enemy.Attack - player.Defense)
	player.Life = max(0, player.Life - enemyDamage)
	enemy.Life = max(0, enemy.Life - playerDamage)
	battle := Battle{
		ID:         uuid.NewString(),
		Enemy:      enemy.Nickname,
		Player:     player.Nickname,
		DiceThrown: diceThrown,
	}
	battles = append(battles, battle)
	json.NewEncoder(w).Encode(battle)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func LoadBattles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(battles)
}
