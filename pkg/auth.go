package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtKey = []byte("secret_key")

// Claims хранит claims в JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AccountData хранит данные пользователей
type AccountData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {

}

// HandleLogin Метод для авторизации пользователя
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var accData AccountData
	// Считывает тело запроса и декодирует их в AccountData
	err := json.NewDecoder(r.Body).Decode(&accData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(os.Getenv("USERNAME"))
	fmt.Println(os.Getenv("PASSWORD"))
	// Сверка данных пользователя с переменными окружения
	if accData.Username != os.Getenv("USERNAME") ||
		accData.Password != os.Getenv("PASSWORD") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Генерация JWT-токена
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: accData.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Отправка JWT-токена клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// JwtAuth Middleware для проверки JWT токена
func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
	})
}
