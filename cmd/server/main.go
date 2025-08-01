package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sp-va/cryptoCurrencies/internal/api"
	"github.com/sp-va/cryptoCurrencies/internal/repository"
	"github.com/sp-va/cryptoCurrencies/internal/service"
	sheduledtasks "github.com/sp-va/cryptoCurrencies/internal/sheduled-tasks"
	"github.com/sp-va/cryptoCurrencies/internal/utils"
)

func init() {
	utils.InitEnvVars()
}

func main() {
	repo, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sheduledtasks.StartCurrencyFetcher(ctx, repo)

	service := service.NewCurrencyService(repo)
	handler := api.NewCurrencyHandler(service)

	router := gin.Default()
	apiV1 := router.Group("/api/v1/currency")
	{
		apiV1.POST("/add", handler.InsertCurrencyHandler)
		apiV1.DELETE("/remove", handler.DeleteCurrencyFromTrackHandler)
		apiV1.GET("/price", handler.GetCoinValue)
	}

	router.Run()
}
