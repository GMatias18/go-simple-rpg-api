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
	ItemID   string `json:"item_id"`
}

type PlayerResponse struct {
	Message string `json:"message"`
}

type Enemy struct {
	Nickname string `json:"nickname"`
	Life     int    `json:"life"`
	Attack   int    `json:"attack"`
	Defense  int    `json:"defense"`
	ItemID   string `json:"item_id"`
}

type Item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	EffectType  string `json:"effect_type"`
	EffectValue int    `json:"effect_value"`
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
var items []Item

func main() {
	initializeItems()

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

	mux.HandleFunc("/item", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddItem(w, r)
		case http.MethodGet:
			LoadItems(w, r)
		}
	})

	fmt.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", mux)
}

func initializeItems() {
	items = []Item{
		{ID: uuid.NewString(), Name: "Espada do Poder", EffectType: "attack", EffectValue: 5},
		{ID: uuid.NewString(), Name: "Escudo de Aço", EffectType: "defense", EffectValue: 3},
		{ID: uuid.NewString(), Name: "Amuleto da Vida", EffectType: "life", EffectValue: 10},
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	item.ID = uuid.NewString()
	items = append(items, item)
	json.NewEncoder(w).Encode(item)
}

func LoadItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

func ApplyItemEffects(player *PlayerRequest, enemy *Enemy) {
	for _, item := range items {
		if player.ItemID == item.ID {
			if item.EffectType == "attack" {
				player.Attack += item.EffectValue
			} else if item.EffectType == "defense" {
				player.Defense += item.EffectValue
			} else if item.EffectType == "life" {
				player.Life += item.EffectValue
			}
		}
		if enemy.ItemID == item.ID {
			if item.EffectType == "attack" {
				enemy.Attack += item.EffectValue
			} else if item.EffectType == "defense" {
				enemy.Defense += item.EffectValue
			} else if item.EffectType == "life" {
				enemy.Life += item.EffectValue
			}
		}
	}
}

func CreateBattle(w http.ResponseWriter, r *http.Request) {
	var battleRequest struct {
		Enemy  string `json:"enemy"`
		Player string `json:"player"`
	}
	if err := json.NewDecoder(r.Body).Decode(&battleRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		http.Error(w, "Player or Enemy not found", http.StatusNotFound)
		return
	}
	ApplyItemEffects(player, enemy)  // Apply item effects before the battle
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
