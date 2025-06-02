package database

import (
	"regexp"
	"testing"
	"time"

	"apiWithDataBase/tasks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetAllTasksDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testDeadline := time.Date(2026, 03, 04, 05, 05, 05, 05, time.UTC)
	testCreatedAt := time.Date(2025, 03, 04, 05, 05, 05, 05, time.UTC)
	rows := mock.NewRows([]string{"id", "title", "description", "status", "priority", "deadline", "created_at"}).AddRow(1, "test", "testDesc", "testStat", "testPrior", testDeadline, testCreatedAt)
	mock.ExpectQuery("SELECT id, title, description, status, priority, deadline, created_at FROM tasks").WillReturnRows(rows)

	getTestTask, err := GetAllTasksDB(db)
	require.NoError(t, err)
	var tasksMass []tasks.Task = []tasks.Task{tasks.Task{ID: 1, Title: "test", Description: "testDesc", Status: "testStat", Priority: "testPrior", Deadline: testDeadline, CreatedAt: testCreatedAt}}
	require.Equal(t, tasksMass, getTestTask)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddTaskDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testDeadline := time.Date(2026, 03, 04, 05, 05, 05, 05, time.UTC)
	testCreatedAt := time.Date(2025, 03, 04, 05, 05, 05, 05, time.UTC)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ($1, $2, $3, $4, $5)")).WithArgs("1", "2", "3", "4", testDeadline).WillReturnResult(sqlmock.NewResult(1, 1))
	var task tasks.Task = tasks.Task{
		ID:          1,
		Title:       "1",
		Description: "2",
		Status:      "3",
		Priority:    "4",
		Deadline:    testDeadline,
		CreatedAt:   testCreatedAt,
	}
	errAdd := AddTaskDB(db, task)
	require.NoError(t, errAdd)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTaskDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM tasks WHERE id = $1")).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	errDelete := DeleteTaskDB(db, 1)
	require.NoError(t, errDelete)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPatchTaskDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE tasks SET title = $1 WHERE id = $2")).WithArgs("matt", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	errPatch := PatchTaskDB(db, "UPDATE tasks SET title = $1 WHERE id = $2", []any{"matt", 1})
	require.NoError(t, errPatch)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertHashedPasswordMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	rows := mock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id")).WithArgs("matt murdock", []byte("john walker")).WillReturnRows(rows)
	id, errHash := InsertHashedPassword(db, "matt murdock", []byte("john walker"))
	require.NoError(t, errHash)
	expectedID := 1
	require.Equal(t, id, expectedID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUserMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	rows := mock.NewRows([]string{"id", "password"}).AddRow(1, "12345")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, password FROM users WHERE username = $1")).WithArgs("hello").WillReturnRows(rows)
	id, password, errCheck := CheckUser(db, "hello")
	require.NoError(t, errCheck)
	expectedPassword := "12345"
	expectedID := 1
	require.Equal(t, expectedID, id)
	require.Equal(t, expectedPassword, password)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	rows := mock.NewRows([]string{"username", "password"}).AddRow("username", "password")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT username, password FROM users WHERE username = $1 LIMIT 1")).WithArgs("username").WillReturnRows(rows)
	user, errGet := GetUserByUsername(db, "username")
	require.NoError(t, errGet)
	var userExp tasks.GetUser = tasks.GetUser{
		Username: "username",
		Password: "password",
	}
	require.Equal(t, user, userExp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserWithTokenInfoMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	date := time.Date(2025, 03, 04, 05, 05, 05, 05, time.UTC)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO userswithtoken (id, username, token_id, expired_at) VALUES ($1, $2, $3, $4)")).WithArgs(1, "matt murdock", "tokenID", date).WillReturnResult(sqlmock.NewResult(1, 1))
	var user tasks.UserWithTokenInfo = tasks.UserWithTokenInfo{
		ID:        1,
		Username:  "matt murdock",
		TokenID:   "tokenID",
		ExpiredAt: date,
	}
	errAdd := AddUserWithTokenInfo(db, user)
	require.NoError(t, errAdd)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckTokenIDMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	date := time.Date(2025, 03, 03, 03, 03, 03, 03, time.UTC)
	rows := mock.NewRows([]string{"username", "expired_at"}).AddRow("me", date)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT username, expired_at FROM userswithtoken WHERE token_id = $1")).WithArgs("privet").WillReturnRows(rows)
	username, getDate, errCheck := CheckTokenID(db, "privet")
	require.NoError(t, errCheck)
	expectedUsername := "me"
	expectedGetDate := time.Date(2025, 03, 03, 03, 03, 03, 03, time.UTC)
	require.Equal(t, username, expectedUsername)
	require.Equal(t, getDate, expectedGetDate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTaskDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	testDeadline := time.Date(2025, 01, 01, 10, 47, 45, 43, time.UTC)
	testCreatedAt := time.Date(2023, 01, 01, 10, 47, 45, 43, time.UTC)
	rows := mock.NewRows([]string{"id", "title", "description", "status", "priority", "deadline", "created_at"}).AddRow(1, "1", "2", "3", "4", testDeadline, testCreatedAt)
	mock.ExpectQuery("SELECT id, title, description, status, priority, deadline, created_at FROM tasks WHERE id = \\$1").WithArgs(1).WillReturnRows(rows)
	gotTask, errGotTask := GetTaskDB(db, 1)
	require.NoError(t, errGotTask)
	var task tasks.Task = tasks.Task{
		ID:          1,
		Title:       "1",
		Description: "2",
		Status:      "3",
		Priority:    "4",
		Deadline:    testDeadline,
		CreatedAt:   testCreatedAt,
	}
	require.Equal(t, gotTask, task)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTasksOverdueDBMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	deadline := time.Date(2025, 03, 04, 03, 03, 03, 03, time.UTC)
	createdAt := time.Date(2024, 03, 04, 03, 03, 03, 03, time.UTC)
	rows := mock.NewRows([]string{"id", "title", "description", "status", "priority", "deadline", "createdAt"}).AddRow(1, "1", "2", "3", "4", deadline, createdAt)
	mock.ExpectQuery(`SELECT id, title, description, status, priority, deadline, created_at FROM tasks WHERE deadline < NOW\(\) ORDER BY deadline ASC;`).WillReturnRows(rows)
	getTask, errGet := GetTasksOverdueDB(db)
	require.NoError(t, errGet)
	var expected []tasks.Task = []tasks.Task{tasks.Task{
		ID:          1,
		Title:       "1",
		Description: "2",
		Status:      "3",
		Priority:    "4",
		Deadline:    deadline,
		CreatedAt:   createdAt,
	}}
	require.Equal(t, getTask, expected)
	require.NoError(t, mock.ExpectationsWereMet())
}
