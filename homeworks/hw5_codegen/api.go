package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

// вы можете использовать ApiError в коде, который получается в результате генерации
// считаем что это какая-то общеизвестная структура
type ApiError struct {
	HTTPStatus int
	Err        error
}

func (ae ApiError) Error() string {
	return ae.Err.Error()
}

// ----------------

const (
	statusUser      = 0
	statusModerator = 10
	statusAdmin     = 20
)

type MyApi struct {
	statuses map[string]int
	users    map[string]*User
	nextID   uint64
	mu       *sync.RWMutex
}

func NewMyApi() *MyApi {
	return &MyApi{
		statuses: map[string]int{
			"user":      0,
			"moderator": 10,
			"admin":     20,
		},
		users: map[string]*User{
			"rvasily": &User{
				ID:       42,
				Login:    "rvasily",
				FullName: "Vasily Romanov",
				Status:   statusAdmin,
			},
		},
		nextID: 43,
		mu:     &sync.RWMutex{},
	}
}

type ProfileParams struct {
	Login string `apivalidator:"required"`
}

type CreateParams struct {
	Login  string `apivalidator:"required,min=10"`
	Name   string `apivalidator:"paramname=full_name"`
	Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	Age    int    `apivalidator:"min=0,max=128"`
}

type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Status   int    `json:"status"`
}

type NewUser struct {
	ID uint64 `json:"id"`
}

// apigen:api {"url": "/user/profile", "auth": false}
func (srv *MyApi) Profile(ctx context.Context, in ProfileParams) (*User, error) {

	if in.Login == "bad_user" {
		return nil, fmt.Errorf("bad user")
	}

	srv.mu.RLock()
	user, exist := srv.users[in.Login]
	srv.mu.RUnlock()
	if !exist {
		return nil, ApiError{http.StatusNotFound, fmt.Errorf("user not exist")}
	}

	return user, nil
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error) {
	if in.Login == "bad_username" {
		return nil, fmt.Errorf("bad user")
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	_, exist := srv.users[in.Login]
	if exist {
		return nil, ApiError{http.StatusConflict, fmt.Errorf("user %s exist", in.Login)}
	}

	id := srv.nextID
	srv.nextID++
	srv.users[in.Login] = &User{
		ID:       id,
		Login:    in.Login,
		FullName: in.Name,
		Status:   srv.statuses[in.Status],
	}

	return &NewUser{id}, nil
}

// 2-я часть
// это похожая структура, с теми же методами, но у них другие параметры!
// код, созданный вашим кодогенератором работает с конкретной струткурой, про другие ничего не знает
// поэтому то что рядом есть ещё походая структура с такими же методами его нисколько не смущает

type OtherApi struct {
}

func NewOtherApi() *OtherApi {
	return &OtherApi{}
}

type OtherCreateParams struct {
	Username string `apivalidator:"required,min=3"`
	Name     string `apivalidator:"paramname=account_name"`
	Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
	Level    int    `apivalidator:"min=1,max=50"`
}

type OtherUser struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Level    int    `json:"level"`
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) {
	return &OtherUser{
		ID:       12,
		Login:    in.Username,
		FullName: in.Name,
		Level:    in.Level,
	}, nil
}
