package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mkhorosh/go-app/models"
	"github.com/mkhorosh/go-app/store"
)

type AnimalRouter struct {
	animalStore *store.AnimalStore
	secretKey   string
}

func NewAnimalRouter(animalStore *store.AnimalStore, secretKey string) *AnimalRouter {
	return &AnimalRouter{animalStore: animalStore, secretKey: secretKey}
}

func (a *AnimalRouter) checkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token required"})
			c.Abort()
			return
		}

		// Убираем префикс "Bearer " из заголовка
		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			return []byte(a.secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

func (a *AnimalRouter) SetupRoutes(router *gin.Engine) {
	animalGroup := router.Group("/animals")
	animalGroup.Use(a.checkAuth())

	{
		animalGroup.GET("/", a.GetAnimals)
		animalGroup.POST("/", a.AddAnimal)
		animalGroup.PUT("/:id", a.UpdateAnimal)
		animalGroup.DELETE("/:id", a.DeleteAnimal)
	}
}

func (a *AnimalRouter) GetAnimals(c *gin.Context) {
	animals, err := a.animalStore.GetAnimals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving animals"})
		return
	}
	c.JSON(http.StatusOK, animals)
}

func (a *AnimalRouter) AddAnimal(c *gin.Context) {
	var animal models.Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	err := a.animalStore.AddAnimal(&animal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving animal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal added successfully"})
}

func (a *AnimalRouter) UpdateAnimal(c *gin.Context) {
	id := c.Param("id")

	var updatedAnimal models.Animal
	if err := c.ShouldBindJSON(&updatedAnimal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	err := a.animalStore.UpdateAnimal(id, &updatedAnimal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating animal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal updated successfully"})
}

func (a *AnimalRouter) DeleteAnimal(c *gin.Context) {
	id := c.Param("id")

	err := a.animalStore.DeleteAnimal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting animal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
