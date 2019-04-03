package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

var uploadFormTmpl = []byte(`
<html>
	<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

type ImgResizeTask struct {
	Name string
	MD5  string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	uploadData, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer uploadData.Close()

	fmt.Fprintf(w, "handler.Filename %v\n", handler.Filename)
	fmt.Fprintf(w, "handler.Header %#v\n", handler.Header)

	tmpName := RandStringRunes(32)

	tmpFile := "./images/" + tmpName + ".jpg"
	newFile, err := os.Create(tmpFile)
	if err != nil {
		http.Error(w, "cant open file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hasher := md5.New()
	writtenBytes, err := io.Copy(newFile, io.TeeReader(uploadData, hasher))
	if err != nil {
		http.Error(w, "cant save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	newFile.Sync()
	newFile.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFile := "./images/" + md5Sum + ".jpg"
	err = os.Rename(tmpFile, realFile)
	if err != nil {
		http.Error(w, "cant raname file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(ImgResizeTask{handler.Filename, md5Sum})

	fmt.Println("put task ", string(data))

	err = rabbitChan.Publish(
		"",                   // exchange
		ImageResizeQueueName, // routing key
		false,                // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	panicOnError("cant publish task", err)

	fmt.Fprintf(w, "Upload %d bytes successful\n", writtenBytes)
}

const (
	ImageResizeQueueName = "image_resize"
)

var (
	rabbitAddr = flag.String("addr", "amqp://guest:guest@192.168.99.100:32778/", "rabbit addr")

	rabbitConn *amqp.Connection

	rabbitChan *amqp.Channel
)

func main() {
	flag.Parse()
	var err error
	rabbitConn, err = amqp.Dial(*rabbitAddr)
	panicOnError("cant connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	panicOnError("cant open chan", err)
	defer rabbitChan.Close()

	q, err := rabbitChan.QueueDeclare(
		ImageResizeQueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	panicOnError("cant init queue", err)

	fmt.Printf("queue %s have %d msg and %d consumers\n",
		q.Name, q.Messages, q.Consumers)

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/upload", uploadPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// никогда так не делайте!
func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + " " + err.Error())
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
