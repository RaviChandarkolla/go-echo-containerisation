package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Score struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}

var rdb = redis.NewClient(&redis.Options{
	Addr: "redis:6379",
})

func incrementScore(w http.ResponseWriter, r *http.Request) {
	var score Score
	json.NewDecoder(r.Body).Decode(&score)

	ctx := context.Background()

	// Atomic increment
	newScore, err := rdb.IncrBy(ctx, fmt.Sprintf("score:%s", score.UserID), 10).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store timestamp
	rdb.Set(ctx, fmt.Sprintf("user:%s:last_update", score.UserID), time.Now().Format(time.RFC3339), 0)

	json.NewEncoder(w).Encode(Score{UserID: score.UserID, Score: int(newScore)})
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Top 10 scores using ZRANGE
	scores, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type LeaderboardEntry struct {
		UserID string  `json:"user_id"`
		Score  float64 `json:"score"`
	}

	var result []LeaderboardEntry
	for _, score := range scores {
		result = append(result, LeaderboardEntry{
			UserID: score.Member.(string),
			Score:  score.Score,
		})
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/score", incrementScore).Methods("POST")
	r.HandleFunc("/leaderboard", getLeaderboard).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
