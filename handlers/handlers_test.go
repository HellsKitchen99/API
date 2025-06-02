package handlers

import (
	"apiWithDataBase/consts"
	"apiWithDataBase/tasks"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func TestDefaultHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	DefaultHandler(response, request)
	expectedBody := "Привет, гость!"

	if expectedBody != response.Body.String() {
		t.Errorf("ожидалось %v, а пришло %v", expectedBody, response.Body.String())
	}

	if response.Code != 200 {
		t.Errorf("ожидалось 200, а пришло %v", response.Code)
	}
}

func setUpTestDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=testapishka sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("не удалось пингануть БД")
	}
	_, errTruncate := db.Exec(`TRUNCATE TABLE tasks RESTART IDENTITY;`)
	if errTruncate != nil {
		return nil, fmt.Errorf("не удалось очистить таблицу")
	}
	return db, nil
}

func setUpTestDBUsers() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=testapishka sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("не удалось пингануть БД")
	}
	_, errTruncate := db.Exec(`TRUNCATE TABLE users RESTART IDENTITY;`)
	if errTruncate != nil {
		return nil, fmt.Errorf("не удалось очистить таблицу")
	}
	return db, nil
}

func InsertTestTasks(db *sql.DB) error {
	time1 := time.Date(2020, 04, 04, 15, 47, 0, 0, time.UTC)
	time2 := time.Date(2020, 04, 04, 15, 48, 0, 0, time.UTC)
	time3 := time.Date(2020, 04, 04, 15, 49, 0, 0, time.UTC)

	_, err1 := db.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('run', 'morning', 'done', 'high', $1)", time1)
	if err1 != nil {
		return fmt.Errorf("ошибка вставки из 1 попытки")
	}

	_, err2 := db.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('cook', 'morning', 'done', 'medium', $1)", time2)
	if err2 != nil {
		return fmt.Errorf("ошибка вставки из 2 попытки")
	}

	_, err3 := db.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('drink', 'morning', 'done', 'low', $1)", time3)
	if err3 != nil {
		return fmt.Errorf("ошибка вставки из 3 попытки")
	}
	return nil
}

func TestGetAllTasks(t *testing.T) {
	dbTest, err := setUpTestDB()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer dbTest.Close()

	err1 := InsertTestTasks(dbTest)
	if err1 != nil {
		t.Errorf("%v", err1)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	hand := GetAllTasks(dbTest)
	hand(w, r)

	var get []tasks.Task
	var expected1 tasks.Task = tasks.Task{ID: 1, Title: "run", Description: "morning", Status: "done", Priority: "high", Deadline: time.Date(2020, 04, 04, 15, 47, 0, 0, time.UTC)}
	var expected2 tasks.Task = tasks.Task{ID: 2, Title: "cook", Description: "morning", Status: "done", Priority: "medium", Deadline: time.Date(2020, 04, 04, 15, 48, 0, 0, time.UTC)}
	var expected3 tasks.Task = tasks.Task{ID: 3, Title: "drink", Description: "morning", Status: "done", Priority: "low", Deadline: time.Date(2020, 04, 04, 15, 49, 0, 0, time.UTC)}
	var expected []tasks.Task = []tasks.Task{expected1, expected2, expected3}
	fmt.Println(w.Body.String())
	decoder := json.NewDecoder(w.Body)
	errDecode := decoder.Decode(&get)
	if errDecode != nil {
		t.Errorf("ну удалось распарсить данные")
	}

	for i := 0; i < 3; i++ {
		if expected[i].ID != get[i].ID {
			t.Errorf("итерация %v - ожидалось %v, а пришло %v", i, expected[i].ID, get[i].ID)
		}
		if expected[i].Title != get[i].Title {
			t.Errorf("итерация %v - ожидалось %v, а пришло %v", i, expected[i].Title, get[i].Title)
		}
		if expected[i].Description != get[i].Description {
			t.Errorf("итерация %v - ожидалось %v, а пришло %v", i, expected[i].Description, get[i].Description)
		}
		if expected[i].Status != get[i].Status {
			t.Errorf("итерация %v - ожидалось %v, а пришло %v", i, expected[i].Status, get[i].Status)
		}
		if expected[i].Priority != get[i].Priority {
			t.Errorf("итерация %v - ожидалось %v, а пришло %v", i, expected[i].Priority, get[i].Priority)
		}
	}

	if w.Code != 200 {
		t.Errorf("ожидалось %v, а пришло %v", 200, w.Code)
	}
}

func TestAddTask(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("не удалось подключиться к тестовой БД")
	}
	defer dbTest.Close()
	date := time.Date(2020, 7, 25, 0, 0, 0, 0, time.UTC)
	var testTask tasks.TestTask = tasks.TestTask{
		Title:       "bike",
		Description: "cycling",
		Status:      "todo",
		Priority:    "high",
		Deadline:    date,
	}
	massBytes, err := json.Marshal(testTask)
	if err != nil {
		t.Errorf("ошибка при парсинге данных в срез байтов")
	}
	buf := bytes.NewBuffer(massBytes)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/tasks", buf)
	hand := AddTask(dbTest)
	hand(w, r)
	expected := "Квест добавлен"
	if w.Body.String() != expected {
		t.Errorf("ожидалось %v, а пришло %v", expected, w.Body.String())
	}
}

func InsertTestTask(db *sql.DB) error {
	date := time.Date(2022, 9, 12, 0, 0, 0, 0, time.UTC)
	_, err := db.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('daredevil', 'watch daredevil', 'todo', 'high', $1)", date)
	if err != nil {
		return fmt.Errorf("не удалось сделать запрос в БД")
	}
	return nil
}

func TestDeleteTask(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	err := InsertTestTask(dbTest)
	if err != nil {
		t.Errorf("%v", errDBTest)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/tasks", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	hand := DeleteTask(dbTest)
	hand(w, r)
	expected := "квест удален"
	if expected != w.Body.String() {
		t.Errorf("ожидалось %v, а пришло %v", expected, w.Body.String())
	}
}

func TestPatchTask(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()

	err := InsertTestTask(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var task tasks.TaskForPatch = tasks.TaskForPatch{
		Deadline: &date,
	}
	massBytes, err := json.Marshal(task)
	if err != nil {
		t.Errorf("не удалось распарсить данные")
	}
	buf := bytes.NewBuffer(massBytes)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/tasks", buf)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	hand := PatchTask(dbTest)
	hand(w, r)
	expected := "квест обновлен"
	body := w.Body.String()

	if expected != body {
		t.Errorf("ожидалось %v, а пришло %v", expected, body)
	}
}

func TestGetTask(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()

	err := InsertTestTasks(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/task", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "3")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	hand := GetTask(dbTest)
	hand(w, r)
	date := time.Date(2020, 04, 04, 15, 49, 0, 0, time.UTC)
	expected := tasks.TestTask{
		ID:          3,
		Title:       "drink",
		Description: "morning",
		Status:      "done",
		Priority:    "low",
		Deadline:    date,
	}
	var outputTestTask tasks.Task
	decoder := json.NewDecoder(w.Body)
	errDecode := decoder.Decode(&outputTestTask)
	if errDecode != nil {
		t.Errorf("не удалось распарсить данные - %v", errDecode)
	}
	if expected.ID != outputTestTask.ID {
		t.Errorf("ожидалось %v, а пришло %v", expected.ID, outputTestTask.ID)
	}
	if expected.Title != outputTestTask.Title {
		t.Errorf("ожидалось %v, а пришло %v", expected.Title, outputTestTask.Title)
	}
	if expected.Description != outputTestTask.Description {
		t.Errorf("ожидалось %v, а пришло %v", expected.Description, outputTestTask.Description)
	}
	if expected.Status != outputTestTask.Status {
		t.Errorf("ожидалось %v, а пришло %v", expected.Status, outputTestTask.Status)
	}
	if expected.Priority != outputTestTask.Priority {
		t.Errorf("ожидалось %v, а пришло %v", expected.Priority, outputTestTask.Priority)
	}
	if expected.Deadline != outputTestTask.Deadline {
		t.Errorf("ожидалось %v, а пришло %v", expected.Deadline, outputTestTask.Deadline)
	}

	if w.Code != 200 {
		t.Errorf("ожидалось %v, а пришло %v", 200, w.Code)
	}
}

func TestGetTasksOverdue(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("ошибка подключения к БД")
	}
	defer dbTest.Close()

	err := InsertTestTasks(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/tasks/overdue", nil)
	hand := GetTasksOverdue(dbTest)
	hand(w, r)
	time1 := time.Date(2020, 04, 04, 15, 47, 0, 0, time.UTC)
	time2 := time.Date(2020, 04, 04, 15, 48, 0, 0, time.UTC)
	time3 := time.Date(2020, 04, 04, 15, 49, 0, 0, time.UTC)
	var output []tasks.Task
	var expected []tasks.TestTask = []tasks.TestTask{
		{ID: 1, Title: "run", Description: "morning", Status: "done", Priority: "high", Deadline: time1},
		{ID: 2, Title: "cook", Description: "morning", Status: "done", Priority: "medium", Deadline: time2},
		{ID: 3, Title: "drink", Description: "morning", Status: "done", Priority: "low", Deadline: time3},
	}
	decoder := json.NewDecoder(w.Body)
	errDecode := decoder.Decode(&output)
	if errDecode != nil {
		t.Errorf("ошибка при декодировании данных")
	}
	for i := 0; i < 3; i++ {
		if expected[i].ID != output[i].ID {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].ID, output[i].ID)
		}
		if expected[i].Title != output[i].Title {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].Title, output[i].Title)
		}
		if expected[i].Description != output[i].Description {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].Description, output[i].Description)
		}
		if expected[i].Status != output[i].Status {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].Status, output[i].Status)
		}
		if expected[i].Priority != output[i].Priority {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].Priority, output[i].Priority)
		}
		if expected[i].Deadline != output[i].Deadline {
			t.Errorf("ожидалось %v, а пришло %v", expected[i].Deadline, output[i].Deadline)
		}
	}
	if len(output) != 3 {
		t.Errorf("ожидалось %v, а пришло %v", 3, len(output))
	}
	if w.Code != 200 {
		t.Errorf("ожидалоось %v, а пришло %v", 200, w.Code)
	}
}

func TestHandlerRegister(t *testing.T) {
	dbTest, errDBTest := setUpTestDBUsers()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	var user tasks.GetUser = tasks.GetUser{
		Username: "john",
		Password: "12345",
	}
	buf, err := json.Marshal(user)
	if err != nil {
		t.Errorf("ошибка кодирования данных")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(buf))
	hand := HandlerRegister(dbTest)
	hand(w, r)
	expected := "регистрация прошла успешно"
	if expected != w.Body.String() {
		t.Errorf("ожидалось %v, а пришло %v", expected, w.Body.String())
	}
	if w.Code != http.StatusCreated {
		t.Errorf("ожидалось %v, а пришло %v", http.StatusCreated, w.Code)
	}
}

func InsertTestUser(db *sql.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("не удалось захешировать пароль")
	}
	_, errExec := db.Exec("INSERT INTO users (username, password) VALUES ('john walker', $1)", hash)
	if errExec != nil {
		return fmt.Errorf("не удалось сделать запрос в БД - %v", err)
	}
	return nil
}

func TestHandlerLogin(t *testing.T) {
	dbTest, errDBTest := setUpTestDBUsers()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	var user tasks.GetUser = tasks.GetUser{
		Username: "john walker",
		Password: "12345",
	}
	buf, err := json.Marshal(user)
	if err != nil {
		t.Errorf("ошибка кодировки данных")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(buf))
	errInsertTestUser := InsertTestUser(dbTest)
	if errInsertTestUser != nil {
		t.Errorf("%v", errInsertTestUser)
	}
	hand := HandlerLogin(dbTest)
	hand(w, r)
	tokenAccessTest := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	})
	fmt.Println(">> RAW BODY:", w.Body.String())
	var token tasks.TestToken
	errDecode := json.NewDecoder(w.Body).Decode(&token)
	if errDecode != nil {
		t.Errorf("ошибка декодирования токена - %v", errDecode)
	}
	parsedToken, errParsedToken := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод")
		}
		return []byte(consts.SecretAccess), nil
	})
	if errParsedToken != nil {
		t.Errorf("токен невалиден")
	}
	realClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("ошибка получения клеимсов")
	}
	testClaims, ok := tokenAccessTest.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("ошибка получения клеимсов")
	}
	realUsername := realClaims["username"].(string)
	testUsername := testClaims["username"].(string)
	if realUsername != testUsername {
		t.Errorf("имена не совпадают, ожидалось %v, а пришло %v", testUsername, realUsername)
	}
}

func InsertTestForTestRefreshAccessToken(db *sql.DB) error {
	username := "walker"
	tokenId := 2
	expiredAt := time.Now().Add(time.Hour)
	_, err := db.Exec("INSERT INTO usersWithToken (username, token_id, expired_at) VALUES ($1, $2, $3)", username, tokenId, expiredAt)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func CheckTestForTestRefreshAccessToken(db *sql.DB, tokenID string) (string, time.Time, error) {
	var username string
	var expiredAt time.Time
	res := db.QueryRow("SELECT username, expired_at FROM usersWithToken WHERE token_id = $1", tokenID)
	err := res.Scan(&username, &expiredAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка получения токен айди из БД - %v", err)
	}
	return username, expiredAt, nil
}

func TestRefreshAccessToken(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	errLoad := godotenv.Load()
	if errLoad != nil {
		t.Errorf("ошибка загрузки файла окружения")
	}
	secretRefresh := os.Getenv(consts.SecretRefresh)
	if secretRefresh == "" {
		t.Errorf("ошибка получения секрета для рефреша")
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/refresh", nil)
	var user tasks.GetUser = tasks.GetUser{
		Username: "walker",
		Password: "12345",
	}
	id := "2"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userID":   1,
		"pairID":   id,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	refreshToken, err := token.SignedString([]byte(secretRefresh))
	if err != nil {
		t.Errorf("ошибка подписи рефреш токена")
	}
	r.AddCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/refresh",
		SameSite: http.SameSiteStrictMode,
	})
	errInsertTestForTestRefreshAccessToken := InsertTestForTestRefreshAccessToken(dbTest)
	if errInsertTestForTestRefreshAccessToken != nil {
		t.Errorf("ошибка вставки теста - %v", errInsertTestForTestRefreshAccessToken)
	}
	hand := RefreshAccessToken(dbTest)
	hand(w, r)
	newTestToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   1,
		"pairID":   id,
		"username": user.Username,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
	})
	secretAccess := os.Getenv(consts.SecretAccess)
	if secretAccess == "" {
		t.Errorf("пустой секрет для аксес токена")
	}
	newAcessToken, errSign := newTestToken.SignedString([]byte(secretAccess))
	if errSign != nil {
		t.Errorf("ошибка подписи нового аксеса")
	}
	var expected tasks.TestToken = tasks.TestToken{
		AccessToken: newAcessToken,
	}
	var output tasks.TestToken
	errDecode := json.NewDecoder(w.Body).Decode(&output)
	if errDecode != nil {
		t.Errorf("ошибка декодирования")
	}
	if expected != output {
		t.Errorf("ожидалось %v, а пришло %v", expected, output)
	}
}

func InsertTestUser2(db *sql.DB) error {
	username := "john"
	password := "12345"
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return fmt.Errorf("ошибка вставки - %v", err)
	}
	return nil
}

func TestGetMe(t *testing.T) {
	db, err := setUpTestDBUsers()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/me", nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "john",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	errLoad := godotenv.Load()
	if errLoad != nil {
		t.Errorf("ошибка загрузки файла окружения")
	}
	tokenString := os.Getenv(consts.SecretAccess)
	if tokenString == "" {
		t.Errorf("ошибка получения секрета для аксеса")
	}
	accessToken, errAccessToken := token.SignedString([]byte(tokenString))
	if errAccessToken != nil {
		t.Errorf("ошибка подписи токена - %v", errAccessToken)
	}
	r.Header.Set("Authorization", "Bearer "+accessToken)
	errTestInsert := InsertTestUser2(db)
	if errTestInsert != nil {
		t.Errorf("%v", errTestInsert)
	}
	var user tasks.GetUser
	hand := GetMe(db)
	hand(w, r)
	fmt.Println(w.Body)
	errDecode := json.NewDecoder(w.Body).Decode(&user)
	if errDecode != nil {
		t.Errorf("ошибка декодирования данных - %v", errDecode)
	}
	var expected tasks.GetUser = tasks.GetUser{
		Username: "john",
		Password: "12345",
	}
	if expected != user {
		t.Errorf("ожидалось %v, а пришло %v", expected, user)
	}
}
