package main

import (
    "github.com/gin-gonic/gin"
    
    "github.com/HarshithRajesh/zapster/initializers"
    "github.com/HarshithRajesh/zapster/controllers"
    "github.com/HarshithRajesh/zapster/middleware" // Ensure this matches the package name
)

func init() {
    initializers.LoadEnvs()
    initializers.ConnectDB()
}

func main() {
    router := gin.Default()
    router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

    router.POST("/auth/signup", controllers.CreateUser )
    router.POST("/auth/login", controllers.Login)
    router.GET("/user/profile", middleware.CheckAuth, controllers.GetUserProfile) // Use middleware instead of middlewares
    router.Run()
}
