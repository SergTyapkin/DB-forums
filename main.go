package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"DB-forums/models"
	"DB-forums/server/handlers"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/prometheus/client_golang/prometheus"
)

var requestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "total_requests",
	})

func middlewareFunc(_ *mux.Router) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			requestsTotal.Inc()
			response.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(response, request)
		})
	}
}

func main() {
	prometheus.MustRegister(requestsTotal)

	parsedConnection, err := pgx.ParseConnectionString("host=localhost user=postgres password=root dbname=DB-forums sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	parsedConnection.PreferSimpleProtocol = true

	connectionConfig := pgx.ConnPoolConfig{
		ConnConfig:     parsedConnection,
		MaxConnections: 69,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	models.DB, err = pgx.NewConnPool(connectionConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/*---------TEST--------
	handlers.TestDropDB()
	handlers.TestInitDB()
	//---------TEST--------*/

	router := mux.NewRouter()

	router.Use(middlewareFunc(router))

	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/create", handlers.ForumCreate).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/details", handlers.ForumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/create", handlers.ThreadCreate).Methods(http.MethodPost)
	router.HandleFunc("/api/forum/{slug}/users", handlers.ForumUsers).Methods(http.MethodGet)
	router.HandleFunc("/api/forum/{slug}/threads", handlers.ForumThreads).Methods(http.MethodGet)

	router.HandleFunc("/api/post/{id:[0-9]+}/details", handlers.PostDetails).Methods(http.MethodGet, http.MethodPost)

	//router.HandleFunc("/api/service/init", handlers.InitDB)
	//router.HandleFunc("/api/service/drop", handlers.DropDB)
	router.HandleFunc("/api/service/clear", handlers.ClearDB).Methods(http.MethodPost)
	router.HandleFunc("/api/service/status", handlers.StatusDB).Methods(http.MethodGet)

	router.HandleFunc("/api/thread/{slug_or_id}/create", handlers.PostsCreate).Methods(http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/details", handlers.ThreadDetails).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/api/thread/{slug_or_id}/posts", handlers.ThreadPosts).Methods(http.MethodGet)
	router.HandleFunc("/api/thread/{slug_or_id}/vote", handlers.VoteCreate).Methods(http.MethodPost)

	router.HandleFunc("/api/user/{nickname}/create", handlers.UserCreate).Methods(http.MethodPost)
	router.HandleFunc("/api/user/{nickname}/profile", handlers.UserProfile).Methods(http.MethodGet, http.MethodPost)

	http.Handle("/", router)

	port := "5000"
	fmt.Println("Server listen at: ", port)
	err = http.ListenAndServe(":"+port, nil)
}
