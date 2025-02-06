package main

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/controller/image"
	"gateway/pkg/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	conf, err := config.Init()

	if err != nil {
		log.Fatalf("error while initializing config: %v", err)
	}

	l, err := logger.NewLogger(
		[]string{conf.KafkaAddress},
		conf.KafkaTopic,
		"gateway",
		100,
	)

	if err != nil {
		log.Fatalf("error while creating logger: %v", err)
	}

	l.Info("successfully initialized logger")
	l.Info("starting gateway service")

	imageCache, err := image.NewCache(conf.ImageCacheDir, conf.ImageBaseURL, &http.Client{}, l)

	if err != nil {
		l.Error(fmt.Sprintf("error while creating image cache: %v", err))
		panic(fmt.Sprintf("error while creating image cache: %v", err))
	}

	l.Info("successfully initialized image cache")

	r := chi.NewRouter()

	l.Info("successfully initialized router")
	l.Info("starting initializing routes")

	r.Get("/img/{filename}", imageCache.GetImage)

	l.Info("all routes initialized")

	l.Info(fmt.Sprintf("server starting on %s", conf.ServerAddress))

	if err := http.ListenAndServe(conf.ServerAddress, r); err != nil {
		l.Error(fmt.Sprintf("error while running server: %v", err))
		panic(fmt.Sprintf("error while running server: %v", err))
	}
}
