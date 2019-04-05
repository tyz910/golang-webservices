package main

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	// "../sess-client/client"
	apiClient "coursera/microservices/swagger/sess-client/client"
	auth "coursera/microservices/swagger/sess-client/client/auth_checker"
	models "coursera/microservices/swagger/sess-client/models"
)

func main() {

	transport := httptransport.New("127.0.0.1:8080", "", []string{"http"})

	client := apiClient.New(transport, strfmt.Default)
	sessManager := client.AuthChecker

	// создаем сессию
	sessId, err := sessManager.Create(auth.NewCreateParams().WithBody(
		&models.SessionSession{
			Login:     "rvasily",
			Useragent: "chrome",
		},
	))
	fmt.Println("sessId", sessId, err)

	// проверяем сессию
	sess, err := sessManager.Check(auth.
		NewCheckParams().
		WithID(sessId.Payload.ID))
	fmt.Println("after create", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(auth.NewDeleteParams().WithBody(
		&models.SessionSessionID{
			ID: sessId.Payload.ID,
		},
	))

	// проверяем еще раз
	sess, err = sessManager.Check(auth.
		NewCheckParams().
		WithID(sessId.Payload.ID))
	fmt.Println("after delete", sess, err)
}
