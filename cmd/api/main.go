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

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" || playerRequest.Life == 0 || playerRequest.Attack == 0 {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname, life and attack is required"})
		return
	}

	if playerRequest.Attack > 10 || playerRequest.Attack <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player attack must be between 1 and 10"})
		return
	}

	if playerRequest.Life > 100 || playerRequest.Life <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player life must be between 1 and 100"})
		return
	}

	for _, player := range players {
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname already exits"})
			return
		}
	}

	player := PlayerRequest{
		Nickname: playerRequest.Nickname,
		Life:     playerRequest.Life,
		Attack:   playerRequest.Attack}
	players = append(players, player)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func LoadPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	for i, player := range players {
		if player.Nickname == nickname {
			players = append(players[:i], players[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func LoadPlayerByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	for _, player := range players {
		if player.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(player)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func SavePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := r.URL.Query().Get("nickname")

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname is required"})
		return
	}

	indexPlayer := -1
	for i, player := range players {
		if player.Nickname == nickname {
			indexPlayer = i
		}
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname already exits"})
			return
		}
	}

	if indexPlayer != -1 {
		players[indexPlayer].Nickname = playerRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(players[indexPlayer])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func AddEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var enemyRequest struct {
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if enemyRequest.Nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname is required"})
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
	newEnemy := Enemy{
		Nickname: enemyRequest.Nickname,
		Life:     rand.Intn(10) + 1,
		Attack:   rand.Intn(10) + 1,
	}

	enemies = append(enemies, newEnemy)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newEnemy)
}

func LoadEnemies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemies)
}

func LoadEnemyByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	for _, enemy := range enemies {
		if enemy.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(enemy)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func UpdateEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	var enemyRequest struct {
		Nickname string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if enemyRequest.Nickname == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname is required"})
		return
	}

	indexEnemy := -1
	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			indexEnemy = i
		}
		if enemy.Nickname == enemyRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname already exists"})
			return
		}
	}

	if indexEnemy != -1 {
		enemies[indexEnemy].Nickname = enemyRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(enemies[indexEnemy])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func DeleteEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	var player PlayerRequest
	var enemy Enemy
	playerFound := false
	enemyFound := false

	for _, p := range players {
		if p.Nickname == battleRequest.Player {
			player = p
			playerFound = true
			break
		}
	}

	for _, e := range enemies {
		if e.Nickname == battleRequest.Enemy {
			enemy = e
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
	damage := player.Attack + diceThrown

	if damage >= enemy.Life {
		enemy.Life = 0
	} else {
		enemy.Life -= damage
	}

	if enemy.Attack > player.Life {
		player.Life = 0
	} else {
		player.Life -= enemy.Attack
	}

	battle := Battle{
		ID:         uuid.New().String(),
		Enemy:      enemy.Nickname,
		Player:     player.Nickname,
		DiceThrown: diceThrown,
	}
	battles = append(battles, battle)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battle)
}

func main() {
	http.HandleFunc("/player", AddPlayer)
	http.HandleFunc("/players", LoadPlayers)
	http.HandleFunc("/player/delete", DeletePlayer)
	http.HandleFunc("/player/update", SavePlayer)
	http.HandleFunc("/player/load", LoadPlayerByNickname)

	http.HandleFunc("/enemy", AddEnemy)
	http.HandleFunc("/enemies", LoadEnemies)
	http.HandleFunc("/enemy/delete", DeleteEnemy)
	http.HandleFunc("/enemy/update", UpdateEnemy)
	http.HandleFunc("/enemy/load", LoadEnemyByNickname)

	http.HandleFunc("/battle", CreateBattle)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
