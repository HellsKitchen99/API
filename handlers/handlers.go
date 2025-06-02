package handlers

import (
	"apiWithDataBase/consts"
	"apiWithDataBase/database"
	"apiWithDataBase/tasks"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, "Привет, гость!")
}

func GetAllTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		tasks, errGetAllTasks := database.GetAllTasksDB(db)
		if errGetAllTasks != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, errGetAllTasks)
			return
		}
		encoder := json.NewEncoder(w)
		errEncoder := encoder.Encode(tasks)
		if errEncoder != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "ошибка парсинга данных в JSON")
			return
		}
	}
}

func AddTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var getTask tasks.Task
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&getTask)
		if errDecode != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Введите верный JSON")
			return
		}
		errAdd := database.AddTaskDB(db, getTask)
		if errAdd != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Введите корректные данные; ошибка - %v", errAdd)
			return
		}
		fmt.Fprint(w, "Квест добавлен")
	}
}

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		idGet := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idGet)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Введите верный ID")
			return
		}
		errDelete := database.DeleteTaskDB(db, id)
		if errDelete != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось подключиться к вашим данным")
			return
		}
		fmt.Fprint(w, "квест удален")
	}
}

func PatchTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t tasks.TaskForPatch
		var args []any
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&t)
		if errDecode != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Введите корректные данные")
			return
		}
		sqlRequest := "UPDATE tasks SET "
		i := 1
		if t.Title != nil {
			sqlRequest += fmt.Sprintf("title = $%v,", i)
			args = append(args, t.Title)
			i++
		}
		if t.Description != nil {
			sqlRequest += fmt.Sprintf("description = $%v,", i)
			args = append(args, t.Description)
			i++
		}
		if t.Status != nil {
			sqlRequest += fmt.Sprintf("status = $%v,", i)
			args = append(args, t.Status)
			i++
		}
		if t.Priority != nil {
			sqlRequest += fmt.Sprintf("priority = $%v,", i)
			args = append(args, t.Priority)
			i++
		}
		if t.Deadline != nil {
			sqlRequest += fmt.Sprintf("deadline = $%v,", i)
			args = append(args, t.Deadline)
			i++
		}

		idGet := chi.URLParam(r, "id")
		id, errParseID := strconv.Atoi(idGet)
		if errParseID != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Введите корректный ID")
			return
		}

		sqlRequestFinal := strings.TrimSuffix(sqlRequest, ",")
		sqlRequestFinal += fmt.Sprintf(" WHERE id = $%v", i)
		args = append(args, id)
		if len(args) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Нечего изменять")
			return
		}

		err := database.PatchTaskDB(db, sqlRequestFinal, args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось подключиться к вашим квестам")
			return
		}
		fmt.Fprint(w, "квест обновлен")
	}
}

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		idGet := chi.URLParam(r, "id")
		id, errParseID := strconv.Atoi(idGet)
		if errParseID != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Введите корректный ID")
			return
		}
		outputTask, errGetTask := database.GetTaskDB(db, id)
		if errGetTask != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось получить квест по заданному ID")
			return
		}
		encoder := json.NewEncoder(w)
		errParseInJSON := encoder.Encode(outputTask)
		if errParseInJSON != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось кодировать данные для передачи")
			return
		}
	}
}

func GetTasksOverdue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		outputTasks, errOverdue := database.GetTasksOverdueDB(db)
		if errOverdue != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось получить доступ к вашим квестам")
			return
		}

		encoder := json.NewEncoder(w)
		errEncode := encoder.Encode(outputTasks)
		if errEncode != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "не удалось кодировать данные для отправвки")
			return
		}
	}
}

func HandlerRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var getUser tasks.GetUser
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&getUser)
		if errDecode != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "некорректные данные")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(getUser.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "ошибка, не удалось провести регистрацию")
		}
		_, errHash := database.InsertHashedPassword(db, getUser.Username, hash)
		if errHash != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "вы уже зарегестрированы")
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "регистрация прошла успешно")
	}
}

func HandlerLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errLoad := godotenv.Load()
		if errLoad != nil {
			log.Print("ошибка загрузки файла окружения")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var getUser tasks.GetUser
		decoder := json.NewDecoder(r.Body)
		errDecode := decoder.Decode(&getUser)
		if errDecode != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "введите корректные данные")
			return
		}
		id, password, err := database.CheckUser(db, getUser.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("ошибка - %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "вы не зарегестрированы")
				return
			}
			log.Printf("ошибка - %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "ошибка сервера")
			return
		}
		errCompare := bcrypt.CompareHashAndPassword([]byte(password), []byte(getUser.Password))
		if errCompare != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "неверный пароль - %v", errCompare)
			return
		}
		tokenId := uuid.New().String()
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":   id,
			"pairID":   tokenId,
			"username": getUser.Username,
			"exp":      time.Now().Add(2 * time.Hour).Unix(),
		})
		secretKeyAccess := os.Getenv(consts.SecretAccess)
		if secretKeyAccess == "" {
			log.Print("secretKeyAccess variable is empty")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		accessTokenString, errAccessTokenString := accessToken.SignedString([]byte(consts.SecretAccess))
		if errAccessTokenString != nil {
			log.Printf("ошибка подписи аксесс токена - %v", errAccessTokenString)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var user tasks.UserWithTokenInfo = tasks.UserWithTokenInfo{
			ID:        id,
			Username:  getUser.Username,
			TokenID:   tokenId,
			ExpiredAt: time.Now().Add(2 * time.Hour),
		}
		errAddUserWithTokenInfo := database.AddUserWithTokenInfo(db, user)
		if errAddUserWithTokenInfo != nil {
			log.Printf("ошибка добавления пользователя с инфо о токене - %v", errAddUserWithTokenInfo)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":   id,
			"pairID":   tokenId,
			"username": getUser.Username,
			"exp":      time.Now().Add(24 * 30 * time.Hour).Unix(),
		})
		secretKeyRefresh := os.Getenv(consts.SecretRefresh)
		if secretKeyRefresh == "" {
			log.Print("secretKeyRefresh variable is empty")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		refreshTokenString, errRefreshTokenString := refreshToken.SignedString([]byte(secretKeyRefresh))
		if errRefreshTokenString != nil {
			log.Printf("ошибка подписи рефреш токена - %v", errRefreshTokenString)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshTokenString,
			HttpOnly: true,
			Secure:   true,
			Path:     "/refresh",
			SameSite: http.SameSiteStrictMode,
		})
		var output tasks.TestToken = tasks.TestToken{
			AccessToken: accessTokenString,
		}
		encoder := json.NewEncoder(w)
		errEncode := encoder.Encode(output)
		if errEncode != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func RefreshAccessToken(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getCookie, errGetCookie := r.Cookie("refresh_token")
		if errGetCookie != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Время сессия истекло, войдите, пожалуйста, заново")
			return
		}
		errErrLoad := godotenv.Load()
		if errErrLoad != nil {
			log.Printf("ошибка загрузки файла окружения - %v", errErrLoad)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		refreshTokenKey := os.Getenv(consts.SecretRefresh)
		if refreshTokenKey == "" {
			log.Print("не удалось получить переменную из файла окружения\n")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		refreshToken := getCookie.Value
		tokenRefresh, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("не тот метод подписи")
			}
			return []byte(refreshTokenKey), nil
		})
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		claims, ok := tokenRefresh.Claims.(jwt.MapClaims)
		if !ok || claims["username"] == nil {
			http.Error(w, "проблемы с токеном", http.StatusUnauthorized)
			return
		}
		tokenId := claims["pairID"].(string)
		usern, expiredAt, errCheckTokenID := database.CheckTokenID(db, tokenId)
		if errCheckTokenID != nil {
			log.Printf("%v", errCheckTokenID)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if time.Now().After(expiredAt) {
			log.Print("токен протух\n")
			http.Error(w, "ваш токен недействителен", http.StatusUnauthorized)
			return
		}
		username := claims["username"]
		if username != usern {
			log.Print("не совпадают имена пользователей\n")
			http.Error(w, "вы не авторизованы", http.StatusUnauthorized)
			return
		}
		exp := claims["exp"].(float64)
		if time.Now().After(time.Unix(int64(exp), 0)) {
			log.Print("токен устарел")
			http.Error(w, "ваша сессия устарела", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["userID"].(float64)
		if !ok {
			log.Print("ошибка приведения userID к верному типу\n")
			http.Error(w, "ошибка с userID", http.StatusUnauthorized)
			return
		}
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":   int(userID),
			"pairID":   tokenId,
			"username": username,
			"exp":      time.Now().Add(2 * time.Hour).Unix(),
		})
		newAccessTokenKey := os.Getenv(consts.SecretAccess)
		if newAccessTokenKey == "" {
			log.Print("ошибка получения переменной из файла окружения\n")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		accessToken, errSigningAccessToken := newAccessToken.SignedString([]byte(newAccessTokenKey))
		if errSigningAccessToken != nil {
			log.Printf("ошибка подписи токена - %v", errSigningAccessToken)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var finalToken tasks.TestToken = tasks.TestToken{
			AccessToken: accessToken,
		}
		encoder := json.NewEncoder(w)
		errEncode := encoder.Encode(finalToken)
		if errEncode != nil {
			log.Printf("ошибка кодировки токена - %v", errEncode)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

func GetMe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Загрузка для теста
		err := godotenv.Load()
		if err != nil {
			fmt.Println("не удалось загрузить файл окружения")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		secretAccess := os.Getenv(consts.SecretAccess)
		if secretAccess == "" {
			fmt.Println("ошибка получения секрета для аксеса")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return []byte(secretAccess), nil
		})
		if err != nil {
			http.Error(w, "вы не авторизованы", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "токен недействителен", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		username := claims["username"].(string)
		user, errGetUser := database.GetUserByUsername(db, username)
		if errGetUser != nil {
			http.Error(w, "вы не авторизованы", http.StatusUnauthorized)
			return
		}
		encoder := json.NewEncoder(w)
		errEncode := encoder.Encode(user)
		if errEncode != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
