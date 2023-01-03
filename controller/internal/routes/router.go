package routes

import (
	"controller-service/internal/handlers"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

func Start(port string) {
	r := chi.NewRouter()
	httpClient := &http.Client{Timeout: time.Second * 5}
	h := &handlers.Handler{HttpClient: httpClient}

	r.Get("/", h.Web)
	r.Post("/", h.Send)

	log.Fatal(http.ListenAndServe(":" + port, r))
}