// Тестирование микросервиса: HTTP сервера и JWT аутентификации (создание фиктивных HTTP запросов и обрабатка HTTP ответов,
// проверяем базовую функциональность каждого эндпоинта, а также правильно ли генерируются JWT токены)
package main

import (
  "bytes"
  "encoding/json"
  "github.com/dgrijalva/jwt-go"
  "net/http"
  "net/http/httptest"
  "testing"
  "time"
)

// TestCreateGame tests creating a new game
func TestCreateGame(t *testing.T) {
  reqBody := bytes.NewBuffer([]byte(`{"players": ["Alice", "Bob"]}`))
  req, err := http.NewRequest("POST", "/create-game", reqBody)
  if err != nil {
    t.Fatalf("Could not create request: %v", err)
  }
  res := httptest.NewRecorder()
  CreateGame(res, req)
  if res.Code != http.StatusOK {
    t.Errorf("Expected status OK; got %v", res.Code)
  }
  var response map[string]string
  if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
    t.Fatalf("Could not decode JSON response: %v", err)
  }
  if _, ok := response["game_id"]; !ok {
    t.Errorf("Expected a game_id in response")
  }
}

// TestGetJWT tests obtaining a JWT token
func TestGetJWT(t *testing.T) {
  game := Game{
    ID:       "1",
    Players:  []string{"Alice", "Bob"},
    IsActive: true,
  }
  mux.Lock()
  games["1"] = game
  mux.Unlock()

  req, err := http.NewRequest("GET", "/get-jwt?username=Alice&game_id=1", nil)
  if err != nil {
    t.Fatalf("Could not create request: %v", err)
  }

  res := httptest.NewRecorder()
  GetJWT(res, req)
  if res.Code != http.StatusOK {
    t.Errorf("Expected status OK; got %v", res.Code)
  }

  var response map[string]string
  if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
    t.Fatalf("Could not decode JSON response: %v", err)
  }
  tokenString, ok := response["token"]
  if !ok {
    t.Errorf("Expected a token in response")
    return
  }

  // Parse the token
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return jwtKey, nil
  })
  if err != nil {
    t.Fatalf("Could not parse token: %v", err)
  }

  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok || !token.Valid {
    t.Fatalf("Invalid token claims")
  }
  if claims["username"].(string) != "Alice" || claims["game_id"].(string) != "1" {
    t.Errorf("Invalid claims data")
  }
}

func TestMain(m *testing.M) {
  http.HandleFunc("/create-game", CreateGame)
  http.HandleFunc("/get-jwt", GetJWT)
  m.Run()
}