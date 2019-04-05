package main

import (
	"fmt"
)

func main() {

	// создаем сессию
	sessId, err := AuthCreateSession(
		&Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	fmt.Println("sessId", sessId, err)

	// проеряем сессию
	sess := AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess)

	// удаляем сессию
	AuthSessionDelete(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess)

}
