package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"os"

	"./models"
	"./redis"
	"./strategies"
)

// Strategy is an interface for pub/sub strategies
type Strategy interface {
	Setup()
	Add(s *models.Connection, channel string) error
	Remove(uuid, channel string) error
	Publish(channel string, r io.Reader) error
}

type server struct {
	Strategy
	log     models.Logger
	Address string
}

// subscribeHandler is the HTTP handler to subscribe a user to a channel
func (s *server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	subscription := models.NewConnection()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Expires", "0")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.

	subscription.Disconnected = w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-subscription.Disconnected
		log.Printf("Client disconnected: %s\n", subscription.ID)
		s.Remove(subscription.ID, channel)
	}()

	if err := s.Add(subscription, channel); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong! %s", err.Error())))
		return
	}

	select {
	case msg := <-subscription.C:
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(msg))
	case <-time.After(30 * time.Second):
		w.WriteHeader(http.StatusNoContent)
	}
	s.Remove(subscription.ID, channel)
}

// publishHanlder is the HTTP handler publish messages to a channel
func (s *server) publishHanlder(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	channel := params.Get("channel")
	s.Publish(channel, r.Body)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Message sent!")))
}

// publishHanlder is the HTTP handler publish messages to a channel
func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong! %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *server) setup() {
	s.Strategy.Setup()
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/subscribe", s.subscribeHandler)
	http.HandleFunc("/publish", s.publishHanlder)
}

func main() {
	var serverAddress = ":8080"
	if len(os.Args) == 2 {
		serverAddress = os.Args[1]
	}

	strategy, err := strategies.NewRedisStrategy(redis.ConfigFromEnv())
	if err != nil {
		panic(err)
	}

	s := server{
		log:      models.StdLogger,
		Strategy: strategy,
		Address:  serverAddress,
	}

	go s.setup()

	log.Printf("server started at %s\n", s.Address)
	log.Fatal(http.ListenAndServe(s.Address, nil))
}
