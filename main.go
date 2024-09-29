package main

// TODO - Frontend JS (1): JavaScript
// TODO - JWT Authorization (2) + http.SetCookie(so-so) + AccessToken and RefreshToken(extra)

// TODO DB: clients info, chat history (3): gorm(postgres) (*)
// TODO: 1.init bd docker
// TODO: 2.use gorm

// TODO - read, write methods: add message author and body, unmarshal binary and marshal to JSON
// TODO - send binary(img, etc.)
// TODO - subdivide code on packages

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"gochat/pkg"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	// Middleware для логирования и обработки ошибок
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Маршрут для авторизации пользователя
	r.Post("/login", pkg.HandleLogin)
	r.With(pkg.JwtAuth).Handle("/chat", pkg.WebSocketHandler{})

	fmt.Println("Server started, waiting for connections on localhost:8000")
	http.ListenAndServe(":8000", r)
}
