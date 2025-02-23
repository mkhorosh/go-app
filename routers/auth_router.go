package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mkhorosh/go-app/models"
	"github.com/mkhorosh/go-app/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthRouter struct {
	userStore *store.FileUserStore
	secretKey string
}

func NewAuthRouter(userStore *store.FileUserStore, secretKey string) *AuthRouter {
	return &AuthRouter{userStore: userStore, secretKey: secretKey}
}

func (a *AuthRouter) SetupRoutes(router *gin.Engine) {
	fmt.Println("Настроены маршруты auth...") // Для отладки
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", a.Login)
		authGroup.POST("/register", a.Register)
	}
}

func (a *AuthRouter) Register(c *gin.Context) {
	fmt.Println("here1")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	fmt.Println("Received user data:", user)

	users, err := a.userStore.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving users"})
		return
	}
	fmt.Println("Existing users:")
	for _, u := range users {
		fmt.Println("User Login:", u.Login, "Name:", u.Name)
	}

	for _, u := range users {
		if u.Login == user.Login {
			c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	user.Password = string(hashedPassword)

	if a.userStore == nil {
		fmt.Println("userStore is nil!")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	err = a.userStore.SaveUser(&user)
	fmt.Println("After calling SaveUser")
	if err != nil {
		fmt.Println("Error saving user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user"})
		return
	}

	fmt.Println("User created successfully:", user)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (a *AuthRouter) Login(c *gin.Context) {
	var loginData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	fmt.Println("sos")

	users, err := a.userStore.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving users"})
		return
	}

	var foundUser *models.User
	for _, user := range users {
		if user.Login == loginData.Login {
			foundUser = user
			break
		}
	}

	if foundUser == nil || bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = foundUser.Login
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": tokenString})
}
