// Автотесты для функции валидации токена
// 1. Создания боя
// 2. Выдача JWT токена 
// 3. Проверка валидного - не валидного токена 
  
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Тест 1 для создания боя
func TestCreateBattle(t *testing.T) {
    data := map[string][]string{"players": {"player1", "player2"}}
    jsonData, err := json.Marshal(data)
    if err != nil {
        t.Fatal("Failed to marshal test data")
    }

    req, err := http.NewRequest("POST", "/createBattle", bytes.NewBuffer(jsonData))
    if err != nil {
        t.Fatal("Failed to create request")
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(createBattle)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("Expected status OK but got %v", rr.Code)
    }

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)
    if _, exists := response["battleID"]; !exists {
        t.Fatal("Expected battleID in response")
    }
}

// Тест 2 для выдачи JWT токена
func TestGetToken(t *testing.T) {
    // Добавим бой в нашу структуру battles для теста
    battleID := "123456"
    battles[battleID] = []string{"player1"}

    data := map[string]string{"battle_id": battleID, "player": "player1"}
    jsonData, err := json.Marshal(data)
    if err != nil {
        t.Fatal("Failed to marshal test data")
    }

    req, err := http.NewRequest("POST", "/getToken", bytes.NewBuffer(jsonData))
    if err != nil {
        t.Fatal("Failed to create request")
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getToken)

    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("Expected status OK but got %v", rr.Code)
    }

    var response map[string]string
    json.Unmarshal(rr.Body.Bytes(), &response)
    if _, exists := response["token"]; !exists {
        t.Fatal("Expected token in response")
    }
}


// Тест 3 для функции валидации токена

func TestValidateToken(t *testing.T) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "player":   "player1",
        "battleID": "1234",
        "exp":      time.Now().Add(time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        t.Fatal("Failed to sign test token")
    }

    player, battleID, err := validateToken(tokenString)
    if err != nil {
        t.Fatal("Failed to validate token")
    }

    if player != "player1" || battleID != "1234" {
        t.Fatalf("Expected player1 and 1234 but got %s and %s", player, battleID)
    }
}

func TestValidateTokenWithInvalidToken(t *testing.T) {
    _, _, err := validateToken("invalid.token.string")
    if err == nil {
        t.Fatal("Expected error for invalid token")
    }
}


