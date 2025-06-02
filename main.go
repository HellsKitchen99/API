package main

import (
	"apiWithDataBase/consts"
	"apiWithDataBase/database"
	"apiWithDataBase/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] - %v", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(db *sql.DB, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "токен не найден", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неверный метод подписи")
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			username := claims["username"].(string)
			_, errGetUser := database.GetUserByUsername(db, username)
			if errGetUser != nil {
				http.Error(w, "Отказано в доступе", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_HOST")
	user := os.Getenv("DB_HOST")
	password := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_HOST")

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	accessSecret := os.Getenv(consts.SecretAccess)
	refreshSecret := os.Getenv(consts.SecretRefresh)
	if accessSecret == "" || refreshSecret == "" {
		fmt.Println("ошибка получения секрета для рефреша и аксеса")
		return
	}

	db, errDB := database.Connect(dns)
	if errDB != nil {
		fmt.Println("не удалось подключиться к БД")
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(logger)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware(db, accessSecret))
		r.Post("/tasks", handlers.AddTask(db))
		r.Delete("/tasks/{id}", handlers.DeleteTask(db))
		r.Patch("/tasks/{id}", handlers.PatchTask(db))
		r.Get("/me", handlers.GetMe(db))
	})

	r.Get("/", handlers.DefaultHandler)
	r.Get("/tasks", handlers.GetAllTasks(db))
	r.Get("/task/{id}", handlers.GetTask(db))
	r.Get("/tasks/overdue", handlers.GetTasksOverdue(db))

	r.Post("/register", handlers.HandlerRegister(db))
	r.Post("/login", handlers.HandlerLogin(db))
	r.Post("/refresh", handlers.RefreshAccessToken(db))
	fmt.Print("сервак слушает на порту :8352\n")

	err := http.ListenAndServe(":8352", r)
	if err != nil {
		fmt.Printf("ошибка при запуске сервака - %v\n", err)
	}
}
