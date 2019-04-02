package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func CheckoutDummy(w http.ResponseWriter, r *http.Request) {

}

var (
	client = &http.Client{Timeout: time.Second}
)

type Case struct {
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
	Path   string
	Query  string
	Auth   bool
	Status int
	Result interface{}
}

const (
	ApiUserCreate  = "/user/create"
	ApiUserProfile = "/user/profile"
)

// CaseResponse
type CR map[string]interface{}

func TestMyApi(t *testing.T) {
	ts := httptest.NewServer(NewMyApi())

	cases := []Case{
		Case{ // успешный запрос
			Path:   ApiUserProfile,
			Query:  "login=rvasily",
			Status: http.StatusOK,
			Result: CR{
				"error": "",
				"response": CR{
					"id":        42,
					"login":     "rvasily",
					"full_name": "Vasily Romanov",
					"status":    20,
				},
			},
		},
		Case{ // успешный запрос - POST
			Path:   ApiUserProfile,
			Method: http.MethodPost,
			Query:  "login=rvasily",
			Status: http.StatusOK,
			Result: CR{
				"error": "",
				"response": CR{
					"id":        42,
					"login":     "rvasily",
					"full_name": "Vasily Romanov",
					"status":    20,
				},
			},
		},
		Case{ // сработала валидация - логин не должен быть пустым
			Path:   ApiUserProfile,
			Query:  "",
			Status: http.StatusBadRequest,
			Result: CR{
				"error": "login must me not empty",
			},
		},
		Case{ // получили ошибку общего назначения - ваш код сам подставил 500
			Path:   ApiUserProfile,
			Query:  "login=bad_user",
			Status: http.StatusInternalServerError,
			Result: CR{
				"error": "bad user",
			},
		},
		Case{ // получили специализированную ошибку - ваш код поставил статус 404 оттуда
			Path:   ApiUserProfile,
			Query:  "login=not_exist_user",
			Status: http.StatusNotFound,
			Result: CR{
				"error": "user not exist",
			},
		},
		// ------
		Case{ // это должен ответить ваш ServeHTTP - если ему пришло что-то неизвестное (например когда он обрабатывает /user/)
			Path:   "/user/unknown",
			Query:  "login=not_exist_user",
			Status: http.StatusNotFound,
			Result: CR{
				"error": "unknown method",
			},
		},
		// ------
		Case{ // создаём юзера
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=mr.moderator&age=32&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusOK,
			Auth:   true,
			Result: CR{
				"error": "",
				"response": CR{
					"id": 43,
				},
			},
		},
		Case{ // юзер действительно создался
			Path:   ApiUserProfile,
			Query:  "login=mr.moderator",
			Status: http.StatusOK,
			Result: CR{
				"error": "",
				"response": CR{
					"id":        43,
					"login":     "mr.moderator",
					"full_name": "Ivan_Ivanov",
					"status":    10,
				},
			},
		},

		Case{ // только POST
			Path:   ApiUserCreate,
			Method: http.MethodGet,
			Query:  "login=mr.moderator&age=32&status=moderator&full_name=GetMethod",
			Status: http.StatusNotAcceptable,
			Auth:   true,
			Result: CR{
				"error": "bad method",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "any_params=123",
			Status: http.StatusForbidden,
			Auth:   false,
			Result: CR{
				"error": "unauthorized",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=mr.moderator&age=32&status=moderator&full_name=New_Ivan",
			Status: http.StatusConflict,
			Auth:   true,
			Result: CR{
				"error": "user mr.moderator exist",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "&age=32&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "login must me not empty",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_m&age=32&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "login len must be >= 10",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_moderator&age=ten&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "age must be int",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_moderator&age=-1&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "age must be >= 0",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_moderator&age=256&status=moderator&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "age must be <= 128",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_moderator&age=32&status=adm&full_name=Ivan_Ivanov",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "status must be one of [user, moderator, admin]",
			},
		},
		Case{ // status по-умолчанию
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=new_moderator3&age=32&full_name=Ivan_Ivanov",
			Status: http.StatusOK,
			Auth:   true,
			Result: CR{
				"error": "",
				"response": CR{
					"id": 44,
				},
			},
		},
		Case{ // обрабатываем неизвестную ошибку
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "login=bad_username&age=32&full_name=Ivan_Ivanov",
			Status: http.StatusInternalServerError,
			Auth:   true,
			Result: CR{
				"error": "bad user",
			},
		},
	}

	runTests(t, ts, cases)
}

func TestOtherApi(t *testing.T) {
	ts := httptest.NewServer(NewOtherApi())

	cases := []Case{
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "username=I3apBap&level=1&class=barbarian&account_name=Vasily",
			Status: http.StatusBadRequest,
			Auth:   true,
			Result: CR{
				"error": "class must be one of [warrior, sorcerer, rouge]",
			},
		},
		Case{
			Path:   ApiUserCreate,
			Method: http.MethodPost,
			Query:  "username=I3apBap&level=1&class=warrior&account_name=Vasily",
			Status: http.StatusOK,
			Auth:   true,
			Result: CR{
				"error": "",
				"response": CR{
					"id":        12,
					"login":     "I3apBap",
					"full_name": "Vasily",
					"level":     1,
				},
			},
		},
	}

	runTests(t, ts, cases)
}

func runTests(t *testing.T, ts *httptest.Server, cases []Case) {
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)

		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		if item.Method == http.MethodPost {
			reqBody := strings.NewReader(item.Query)
			req, err = http.NewRequest(item.Method, ts.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest(item.Method, ts.URL+item.Path+"?"+item.Query, nil)
		}

		if item.Auth {
			req.Header.Add("X-Auth", "100500")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		// fmt.Printf("[%s] body: %s\n", caseName, string(body))

		if resp.StatusCode != item.Status {
			t.Errorf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Errorf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		// reflect.DeepEqual не работает если нам приходят разные типы
		// а там приходят разные типы (string VS interface{}) по сравнению с тем что в ожидаемом результате
		// этот маленький грязный хак конвертит данные сначала в json, а потом обратно в interface - получаем совместимые результаты
		// не используйте это в продакшен-коде - надо явно писать что ожидается интерфейс или использовать другой подход с точным форматом ответа
		data, err := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("[%d] results not match\nGot: %#v\nExpected: %#v", idx, result, item.Result)
			continue
		}
	}
}
