package main

import (
	"context"
	"fmt"

	"github.com/bitly/adapter/handler/http"
	"github.com/bitly/adapter/repository"
	"github.com/bitly/config"
	"github.com/bitly/domain"
	"github.com/bitly/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Finished Loading Up env...")

	appConfig, err := config.NewApp()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db, err := config.NewDB().Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Finished Connecting Database...")

	err = db.AutoMigrate(&domain.URLMapping{}, &domain.URLClick{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Finished Migrating Database...")

	ctx := context.Background()

	mappingRepository := repository.NewMappingRepository(db, ctx)
	mappingService := usecase.NewMappingService(mappingRepository, appConfig.LinkTTL)
	handler := http.NewHandler(*mappingService)

	r := gin.Default()
	handler.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err.Error())
	}
}
