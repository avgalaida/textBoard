package main

import (
	"fmt"
	"github.com/avgalaida/textBoard/db"
	"github.com/avgalaida/textBoard/event"
	"github.com/avgalaida/textBoard/search"
	"github.com/avgalaida/textBoard/util"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"time"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", listPostsHandler).
		Methods(http.MethodGet)
	router.HandleFunc("/search", searchPostsHandler).
		Methods(http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	util.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	util.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := search.NewElastic(fmt.Sprintf("http://%s", cfg.ElasticsearchAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		search.SetRepository(es)
		return nil
	})
	defer search.Close()

	util.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		err = es.OnPostCreated(onPostCreated)
		if err != nil {
			log.Println(err)
			return err
		}
		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
