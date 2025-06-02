package database

import (
	"apiWithDataBase/tasks"
	"database/sql"
	"fmt"
	"time"
)

func GetAllTasksDB(db *sql.DB) ([]tasks.Task, error) {
	var Tasks []tasks.Task
	rows, err := db.Query("SELECT id, title, description, status, priority, deadline, created_at FROM tasks ")
	if err != nil {
		return Tasks, fmt.Errorf("ошибка запроса к БД - %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var t tasks.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Deadline, &t.CreatedAt)
		if err != nil {
			return Tasks, fmt.Errorf("ошибка сканирования полей - %v", err)
		}
		Tasks = append(Tasks, t)
	}
	return Tasks, nil
}

func AddTaskDB(db *sql.DB, getTask tasks.Task) error {
	_, errExec := db.Exec("INSERT INTO tasks (title, description, status, priority, deadline) VALUES ($1, $2, $3, $4, $5)", getTask.Title, getTask.Description, getTask.Status, getTask.Priority, getTask.Deadline)
	if errExec != nil {
		return fmt.Errorf("ошибка запроса в БД - %v", errExec)
	}
	return nil
}

func DeleteTaskDB(db *sql.DB, id int) error {
	res, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка запроса к БД на удаление")
	}
	rowsAffected, errRowsAffected := res.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("ошибка получения количества строк")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ничего не удалено")
	}
	return nil
}

func PatchTaskDB(db *sql.DB, sqlReq string, args []interface{}) error {
	res, err := db.Exec(sqlReq, args...)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД - %v", err)
	}
	rowsAffected, errRowsAffected := res.RowsAffected()
	if errRowsAffected != nil {
		return fmt.Errorf("ошибка получения количества строк")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ничего не изменено")
	}
	return nil
}

func GetTaskDB(db *sql.DB, id int) (tasks.Task, error) {
	var t tasks.Task
	res := db.QueryRow("SELECT id, title, description, status, priority, deadline, created_at FROM tasks WHERE id = $1", id)
	err := res.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.Deadline,
		&t.CreatedAt,
	)

	if err != nil {
		return t, fmt.Errorf("не удалось сделать запрос в БД - %v", err)
	}
	return t, nil
}

func GetTasksOverdueDB(db *sql.DB) ([]tasks.Task, error) {
	var OverdueTasks []tasks.Task
	rows, err := db.Query("SELECT id, title, description, status, priority, deadline, created_at FROM tasks WHERE deadline < NOW() ORDER BY deadline ASC;")
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД - %v", err)
	}
	for rows.Next() {
		var t tasks.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Deadline, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("не удалось сканировать поля из БД")
		}
		OverdueTasks = append(OverdueTasks, t)
	}
	return OverdueTasks, nil
}

func InsertHashedPassword(db *sql.DB, username string, hashedPassword []byte) (int, error) {
	res := db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, hashedPassword)
	var id int
	errScan := res.Scan(&id)
	if errScan != nil {
		return 0, fmt.Errorf("не удалось забрать id - %v", errScan)
	}
	return id, nil
}

func CheckUser(db *sql.DB, username string) (int, string, error) {
	var id int
	var password string
	res := db.QueryRow("SELECT id, password FROM users WHERE username = $1", username)
	err := res.Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("ничего не затронуто - %w", err)
		}
		return 0, "", fmt.Errorf("ошибка сканирования - %v", err)
	}
	return id, password, nil
}

func GetUserByUsername(db *sql.DB, username string) (tasks.GetUser, error) {
	var user tasks.GetUser
	res := db.QueryRow("SELECT username, password FROM users WHERE username = $1 LIMIT 1", username)
	err := res.Scan(&user.Username, &user.Password)
	if err != nil {
		return user, fmt.Errorf("ошибка сканирования")
	}
	return user, nil
}

func AddUserWithTokenInfo(db *sql.DB, user tasks.UserWithTokenInfo) error {
	_, err := db.Exec("INSERT INTO userswithtoken (id, username, token_id, expired_at) VALUES ($1, $2, $3, $4)", user.ID, user.Username, user.TokenID, user.ExpiredAt)
	if err != nil {
		return fmt.Errorf("не удалось сделать запрос в БД на добавление пользователя с инфой о токене - %v", err)
	}
	return nil
}

func CheckTokenID(db *sql.DB, tokenID string) (string, time.Time, error) {
	var username string
	var expiredAt time.Time
	res := db.QueryRow("SELECT username, expired_at FROM userswithtoken WHERE token_id = $1", tokenID)
	err := res.Scan(&username, &expiredAt)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ошибка получения токен айди из БД - %v", err)
	}
	return username, expiredAt, nil
}
