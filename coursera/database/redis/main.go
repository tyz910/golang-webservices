package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	redisAddr = flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")

	sessManager *SessionManager

	users = map[string]string{
		"rvasily":        "test",
		"romanov.vasily": "100500",
	}
)

func checkSession(r *http.Request) (*Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess := sessManager.Check(&SessionID{
		ID: cookieSessionID.Value,
	})
	return sess, nil
}

func innerPage(w http.ResponseWriter, r *http.Request) {
	sess, err := checkSession(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sess == nil {
		w.Write(loginFormTmpl)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome, "+sess.Login+" <br />")
	fmt.Fprintln(w, "Session ua: "+sess.Useragent+" <br />")
	fmt.Fprintln(w, `<a href="/logout">logout</a>`)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue("login")
	inputPass := r.FormValue("password")
	expiration := time.Now().Add(24 * time.Hour)

	pass, exist := users[inputLogin]
	if !exist || pass != inputPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sess, err := sessManager.Create(&Session{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		log.Println("cant create session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	flag.Parse()

	var err error
	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Fatalf("cant connect to redis")
	}

	sessManager = NewSessionManager(redisConn)

	http.HandleFunc("/", innerPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessManager.Delete(&SessionID{
		ID: session.Value,
	})

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

var loginFormTmpl = []byte(`
	<html>
		<body>
		<form action="/login" method="post">
			Login: <input type="text" name="login">
			Password: <input type="password" name="password">
			<input type="submit" value="Login">
		</form>
		</body>
	</html>
	`)
