package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/controllers"
  "github.com/HarshithRajesh/zapster/middleware"
)

func init(){
  initializers.LoadEnvs()
  initializers.ConnectDB()
}

func main() {
  r := gin.Default()

  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"data": "hello world"})    
  })
  r.POST("/auth/signup",controllers.CreateUser)
  r.POST("/auth/login",controllers.Login)
  r.POST("/auth/profile",middleware.CheckAuth,controllers.GetUserProfile)

  r.Run()
}
