package main

import (
	"fmt"
	"net/http"

	"github.com/mkhorosh/go-app/routers"
	"github.com/mkhorosh/go-app/store"

	"github.com/gin-gonic/gin"
)

const port = "8080"

func main() {

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})

	userStore := store.NewFileUserStore("data/users.json")
	userRouter := routers.NewUserRouter(userStore)
	userRouter.SetupRoutes(router)

	authRouter := routers.NewAuthRouter(userStore, "secret_key")
	authRouter.SetupRoutes(router)

	animalStore := store.NewAnimalStore("data/animals.json")
	animalRouter := routers.NewAnimalRouter(animalStore, "secret_key")
	animalRouter.SetupRoutes(router)

	fmt.Println("Server is running on port", port)

	for _, route := range router.Routes() {
		fmt.Printf("%s %s\n", route.Method, route.Path)
	}

	router.Run(":" + port)
}
