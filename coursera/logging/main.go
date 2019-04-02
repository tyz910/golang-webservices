package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

type AccessLogger struct {
	StdLogger    *log.Logger
	ZapLogger    *zap.SugaredLogger
	LogrusLogger *logrus.Entry
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		fmt.Printf("FMT [%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		log.Printf("LOG [%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		ac.StdLogger.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		ac.ZapLogger.Info(r.URL.Path,
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

// -----------

func main() {

	addr := "localhost"
	port := 8080

	// std
	fmt.Printf("STD starting server at %s:%d", addr, port)

	// std
	log.Printf("STD starting server at %s:%d", addr, port)

	// zap
	// у zap-а нет логгера по-умолчанию
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	zapLogger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", addr),
		zap.Int("port", port),
	)

	// logrus
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   addr,
		"port":   port,
	}).Info("Starting server")

	AccessLogOut := new(AccessLogger)

	// std
	AccessLogOut.StdLogger = log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)

	// zap
	sugar := zapLogger.Sugar().With(
		zap.String("mode", "[access_log]"),
		zap.String("logger", "ZAP"),
	)
	AccessLogOut.ZapLogger = sugar

	// logrus
	contextLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AccessLogOut.LogrusLogger = contextLogger

	// server stuff
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", mainPage)
	siteHandler := AccessLogOut.accessLogMiddleware(siteMux)
	http.ListenAndServe(":8080", siteHandler)
}
