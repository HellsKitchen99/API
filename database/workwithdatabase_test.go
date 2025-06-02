package database

import (
	"apiWithDataBase/tasks"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setUpTestDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=testapishka sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("ошибка пинга БД - %v", errPing)
	}
	_, errTruncate := db.Exec(`TRUNCATE TABLE tasks RESTART IDENTITY;`)
	if errTruncate != nil {
		return nil, fmt.Errorf("не удалось очистить таблицу")
	}
	return db, nil
}

func TestGetAllTasksDB(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date1 := time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2022, 4, 4, 0, 0, 0, 0, time.UTC)
	_, err1 := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('d', 'a', 'r', 'e', $1)", date1)
	if err1 != nil {
		t.Errorf("не удалось сделать запрос в БД")
	}
	_, err2 := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('m', 'a', 't', 't', $1)", date2)
	if err2 != nil {
		t.Errorf("не удалось сделать запрос в БД")
	}

	output, errGet := GetAllTasksDB(dbTest)
	if errGet != nil {
		t.Errorf("%v", errGet)
	}
	if len(output) != 2 {
		t.Errorf("ожидалось %v, а пришло %v", 2, len(output))
	}
}

func GetAddedTestTask(db *sql.DB) (tasks.Task, error) {
	res := db.QueryRow("SELECT id, title, description, status, priority, deadline, created_at FROM tasks")
	var out tasks.Task
	err := res.Scan(&out.ID, &out.Title, &out.Description, &out.Status, &out.Priority, &out.Deadline, &out.CreatedAt)
	if err != nil {
		return out, fmt.Errorf("не удалось получить данные из БД - %v", err)
	}
	return out, nil
}

func TestAddTaskDB(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date := time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC)
	var inputTestTask tasks.Task = tasks.Task{
		ID:          1,
		Title:       "darevevil",
		Description: "born again",
		Status:      "muse",
		Priority:    "punisher",
		Deadline:    date,
		CreatedAt:   time.Now(),
	}
	err := AddTaskDB(dbTest, inputTestTask)
	if err != nil {
		t.Errorf("не удалось добавить квест")
	}
	output, errGet := GetAddedTestTask(dbTest)
	if errGet != nil {
		t.Errorf("%v", errGet)
	}
	if output.ID != inputTestTask.ID {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.ID, output.ID)
	}
	if output.Title != inputTestTask.Title {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.Title, output.Title)
	}
	if output.Description != inputTestTask.Description {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.Description, output.Description)
	}
	if output.Status != inputTestTask.Status {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.Status, output.Status)
	}
	if output.Priority != inputTestTask.Priority {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.Priority, output.Priority)
	}
	if output.Deadline.Format("2006-01-02 15:04:05") != inputTestTask.Deadline.Format("2006-01-02 15:04:05") {
		t.Errorf("ожидалось %v, а пришло %v", inputTestTask.Deadline.Format("2006-01-02 15:04:05"), output.Deadline.Format("2006-01-02 15:04:05"))
	}
}

func TestDeleteTaskDB(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date := time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC)
	_, err := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('night', 'gown', 'man', 'red', $1)", date)
	if err != nil {
		t.Errorf("ошибка запроса в БД")
	}
	errDelete1 := DeleteTaskDB(dbTest, 1)
	if errDelete1 != nil {
		t.Errorf("ошибка должна быть nil")
	}
	errDelete2 := DeleteTaskDB(dbTest, 2)
	if errDelete2 == nil {
		t.Errorf("должна быть ошибка - %v", errDelete2)
	}
}

func TestPatchTask(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date := time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC)
	_, err := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('r', 'e', 'd', 'd', $1)", date)
	if err != nil {
		t.Errorf("ошибка запроса в БД")
	}
	errPatchTask1 := PatchTaskDB(dbTest, "UPDATE tasks SET title = $1 WHERE id = $2", []any{"nightgown man", 1})
	if errPatchTask1 != nil {
		t.Errorf("ошибка должна быть nil")
	}
	errPatchTask2 := PatchTaskDB(dbTest, "UPDATE tasks SET title = $1 WHERE id = $2", []any{"nightgown man", 2})
	if errPatchTask2 == nil {
		t.Errorf("ошибка должна быть - %v", errPatchTask2)
	}
}

func TestGetTaskDB(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date := time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC)
	_, err := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('daredevil', 'born', 'again', 'masterpiece', $1)", date)
	if err != nil {
		t.Errorf("ошибка запролса в БД")
	}
	_, errGetTaskDB1 := GetTaskDB(dbTest, 2)
	if errGetTaskDB1 == nil {
		t.Errorf("ошибка должна быть - %v", errGetTaskDB1)
	}
	task2, errGetTaskDB2 := GetTaskDB(dbTest, 1)
	if errGetTaskDB2 != nil {
		t.Errorf("ошибка должна быть nil")
	}
	var Task tasks.TaskForTestingDataBase = tasks.TaskForTestingDataBase{
		Title:       "daredevil",
		Description: "born",
		Status:      "again",
		Priority:    "masterpiece",
		Deadline:    date,
	}
	if Task.Title != task2.Title {
		t.Errorf("ожидалось %v, а пришло %v", Task.Title, task2.Title)
	}
	if Task.Description != task2.Description {
		t.Errorf("ожидалось %v, а пришло %v", Task.Description, task2.Description)
	}
	if Task.Status != task2.Status {
		t.Errorf("ожидалось %v, а пришло %v", Task.Status, task2.Status)
	}
	if Task.Priority != task2.Priority {
		t.Errorf("ожидалось %v, а пришло %v", Task.Priority, task2.Priority)
	}
}

func TestGetTasksOverdueDB(t *testing.T) {
	dbTest, errDBTest := setUpTestDB()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	date1 := time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC)
	res1, err1 := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('matt', 'murdock', 'daredevil', 'born again', $1)", date1)
	if err1 != nil {
		t.Errorf("ошибка 1 запроса")
	}
	rowAffected1, errRowAffected1 := res1.RowsAffected()
	if errRowAffected1 != nil {
		t.Errorf("ошибка получения количества строк")
	}
	if rowAffected1 == 0 {
		t.Errorf("ничего не затронуто")
	}
	date2 := time.Date(2024, 03, 04, 0, 0, 0, 0, time.UTC)
	res2, err2 := dbTest.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ('punisher', 'frank', 'castle', 'no they', $1)", date2)
	if err2 != nil {
		t.Errorf("ошибка запроса в БД")
	}
	rowAffected2, errRowAffected2 := res2.RowsAffected()
	if errRowAffected2 != nil {
		t.Errorf("ошибка получения количества строк")
	}
	if rowAffected2 == 0 {
		t.Errorf("ничего не затронуто")
	}
	output, err := GetTasksOverdueDB(dbTest)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(output) != 2 {
		t.Errorf("ожидалось %v, а получили %v", 2, len(output))
	}
	expected := []tasks.TaskForTestingOverdue{
		{Title: "punisher", Description: "frank", Status: "castle", Priority: "no they"},
		{Title: "matt", Description: "murdock", Status: "daredevil", Priority: "born again"},
	}
	for i := 0; i < 2; i++ {
		if expected[i].Title != output[i].Title {
			t.Errorf("ожидалось %v, а получили %v", expected[i].Title, output[i].Title)
		}
		if expected[i].Description != output[i].Description {
			t.Errorf("ожидалось %v, а получили %v", expected[i].Description, output[i].Description)
		}
		if expected[i].Status != output[i].Status {
			t.Errorf("ожидалось %v, а получили %v", expected[i].Status, output[i].Status)
		}
		if expected[i].Priority != output[i].Priority {
			t.Errorf("ожидалось %v, а получили %v", expected[i].Priority, output[i].Priority)
		}
	}
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

func setUpTestDBUsersWithToken() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password='C529@4H0OdEO%}Y' dbname=testapishka sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД")
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("не удалось пингануть БД")
	}
	_, errTruncate := db.Exec(`TRUNCATE TABLE userswithtoken RESTART IDENTITY;`)
	if errTruncate != nil {
		return nil, fmt.Errorf("не удалось очистить таблицу - %v", errTruncate)
	}
	return db, nil
}

func TestInsertHashedPassword(t *testing.T) {
	db, err := setUpTestDBUsers()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()
	username := "matt murdock"
	hashedPassword := []byte("hellsKitchen")
	id, errTestFunc := InsertHashedPassword(db, username, hashedPassword)
	if errTestFunc != nil {
		t.Errorf("%v", err)
	}
	expectedId := 1
	if expectedId != id {
		t.Errorf("ожидалось %v, а пришло %v", expectedId, id)
	}
	expectedUsername := "matt murdock"
	expectedPassword := "hellsKitchen"
	res := db.QueryRow("SELECT username, password FROM users")
	errScan := res.Scan(&expectedUsername, &expectedPassword)
	if errScan != nil {
		t.Errorf("ошибка сканирования")
	}
	if expectedUsername != username || expectedPassword != string(hashedPassword) {
		t.Errorf("короч тест падает окончательно")
	}
}

func TestCheckUser(t *testing.T) {
	db, err := setUpTestDBUsers()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()
	_, errExec := db.Exec("INSERT INTO users (username, password) VALUES ('john walker', 'qwert')")
	if errExec != nil {
		t.Errorf("ошибка вставки")
	}
	id, password, errCheck := CheckUser(db, "john walker")
	if errCheck != nil {
		t.Errorf("%v", errCheck)
	}
	expectedID := 1
	expectedPassword := "qwert"
	if expectedID != id || expectedPassword != password {
		t.Errorf("тест упал")
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, err := setUpTestDBUsers()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()
	_, errExec := db.Exec("INSERT INTO users (username, password) VALUES ('john walker', 'matt murdock')")
	if errExec != nil {
		t.Errorf("ошабка вставки")
	}
	user, errGetUser := GetUserByUsername(db, "john walker")
	if errGetUser != nil {
		t.Errorf("%v", errGetUser)
	}
	var expectedUser tasks.GetUser = tasks.GetUser{
		Username: "john walker",
		Password: "matt murdock",
	}
	if user.Username != expectedUser.Username {
		t.Errorf("ожидалось %v, а пришло %v", expectedUser.Username, user.Username)
	}
	if user.Password != expectedUser.Password {
		t.Errorf("ожидалось %v, а пришло %v", expectedUser.Password, user.Password)
	}
}

func TestAddUserWithTokenInfo(t *testing.T) {
	dbTest, errDBTest := setUpTestDBUsersWithToken()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	var user tasks.UserWithTokenInfo = tasks.UserWithTokenInfo{
		ID:        1,
		Username:  "beetlejuice",
		TokenID:   "nothing",
		ExpiredAt: time.Now().Add(time.Hour),
	}
	err := AddUserWithTokenInfo(dbTest, user)
	if err != nil {
		t.Errorf("%v", err)
	}
	res := dbTest.QueryRow("SELECT id, username, token_id, expired_at FROM userswithtoken")
	var id int
	var username, tokenid string
	var expiredat time.Time
	errScan := res.Scan(&id, &username, &tokenid, &expiredat)
	if errScan != nil {
		t.Errorf("ошибка сканирования - %v", errScan)
	}
	expectedID := 1
	expectedUsername := "beetlejuice"
	expectedTokenID := "nothing"
	if expectedID != id {
		t.Errorf("ожидалось %v, а пришло %v", expectedID, id)
	}
	if expectedUsername != username {
		t.Errorf("ожидалось %v, а пришло %v", expectedUsername, username)
	}
	if expectedTokenID != tokenid {
		t.Errorf("ожидалось %v, а пришло %v", expectedTokenID, tokenid)
	}
}

func TestCheckTokenID(t *testing.T) {
	dbTest, errDBTest := setUpTestDBUsersWithToken()
	if errDBTest != nil {
		t.Errorf("%v", errDBTest)
	}
	defer dbTest.Close()
	tokenID := "murdock"
	time1 := time.Now().Add(time.Hour)
	_, err := dbTest.Exec("INSERT INTO userswithtoken (username, token_id, expired_at) VALUES ('matt', 'murdock', $1)", time1)
	if err != nil {
		t.Errorf("ошибка вставки")
	}
	username, _, errCheckTokenID := CheckTokenID(dbTest, tokenID)
	if errCheckTokenID != nil {

	}
	expectedUsername := "matt"
	if expectedUsername != username {
		t.Errorf("ожидалось %v, а пришло %v", expectedUsername, username)
	}
}
