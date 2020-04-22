package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhl1402/covidscript/cmd/api/playground"
	"github.com/dhl1402/covidscript/cmd/api/svc"

	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	logger := svc.MakeLogger()
	defer logger.Sync()
	sugar := logger.Sugar()
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://2d7eee162843474d8021b4f3592e3df9@o316001.ingest.sentry.io/5209601",
	})

	s := playground.New(sugar.With("svc", "playground"))
	playground.NewHandler(s, router, sugar)

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, panic interface{}) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("500: Internal Server Error"))
		if p, ok := panic.(string); ok {
			w.Write([]byte("\nPanic: " + p))
		}
	}

	errs := make(chan error, 2)
	go func() {
		port := "8080"
		if p := os.Getenv("PORT"); p != "" {
			port = p
		}
		errs <- http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-c
		// TODO: received exit signal: do graceful shutdown
		errs <- nil
	}()

	<-errs
	// TODO: handle error
}
