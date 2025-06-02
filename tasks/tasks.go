package tasks

import (
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskForPatch struct {
	ID          *int       `json:"id"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
}

type TestTask struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
}

type TaskForTestingDataBase struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Deadline    time.Time `json:"deadline"`
}
type TaskForTestingOverdue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TestToken struct {
	AccessToken string `json:"access_token"`
}

type UserWithTokenInfo struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	TokenID   string    `json:"token_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
