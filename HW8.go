// Микросервис авторизации

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")
var battles = make(map[string][]string)

func createBattle(w http.ResponseWriter, r *http.Request) {
    var participants struct {
        Players []string `json:"players"`
    }

    if err := json.NewDecoder(r.Body).Decode(&participants); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    battleID := fmt.Sprintf("%d", time.Now().Unix())
    battles[battleID] = participants.Players

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"battleID": battleID})
}

func getToken(w http.ResponseWriter, r *http.Request) {
    var request struct {
        BattleID string `json:"battle_id"`
        Player   string `json:"player"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    players, ok := battles[request.BattleID]
    if !ok {
        http.Error(w, "Battle not found", http.StatusNotFound)
        return
    }

    for _, player := range players {
        if player == request.Player {
            token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                "player":   request.Player,
                "battleID": request.BattleID,
                "exp":      time.Now().Add(time.Hour * 24).Unix(),
            })

            tokenString, err := token.SignedString(jwtKey)
            if err != nil {
                http.Error(w, "Failed to generate token", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
            return
        }
    }

    http.Error(w, "Player not part of the battle", http.StatusUnauthorized)
}

func main() {
    http.HandleFunc("/createBattle", createBattle)
    http.HandleFunc("/getToken", getToken)

    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Микросервис проверки JWT

func validateToken(tokenStr string) (string, string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return "", "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        player := claims["player"].(string)
        battleID := claims["battleID"].(string)
        return player, battleID, nil
    }

    return "", "", fmt.Errorf("invalid token")
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    tokenStr := strings.TrimPrefix(authHeader, "Bearer")

    player, battleID, err := validateToken(tokenStr)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Обработка сообщения используя player и battleID

    w.Write([]byte("Message processed successfully"))
}

// In the main function:
http.HandleFunc("/handleMessage", handleMessage)
