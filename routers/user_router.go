package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkhorosh/go-app/store"
)

type UserRouter struct {
	userStore *store.FileUserStore
}

func NewUserRouter(userStore *store.FileUserStore) *UserRouter {
	return &UserRouter{userStore: userStore}
}

func (r *UserRouter) SetupRoutes(router *gin.Engine) {
	router.GET("/users", r.GetUsers)
}

func (r *UserRouter) GetUsers(c *gin.Context) {
	users, err := r.userStore.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read users data"})
		return
	}
	c.JSON(http.StatusOK, users)
}
