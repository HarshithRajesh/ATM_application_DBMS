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
    router.GET("/user/profile", middleware.CheckAuth, controllers.GetUserProfile)
    router.POST("/accounts", controllers.CreateAccount)
	  router.GET("/accounts/:id", controllers.GetAccount)
	  router.PUT("/accounts/:id", controllers.UpdateAccount)
	  router.DELETE("/accounts/:id", controllers.DeleteAccount)
	  router.GET("/accounts", controllers.ListAccounts)
    router.POST("/create-card",controllers.CreateCard)
    router.POST("/update-card-status", controllers.UpdateCardStatus) // Update Card Status
    router.POST("/lock-card", controllers.LockCardAfterFailedAttempts) // Lock Card After Failed Attempts
    router.DELETE("/delete-card/:card_number", controllers.DeleteCard) // Delete Card
    router.POST("/transactions/withdraw", controllers.CashWithdrawal)
    router.POST("/transactions/deposit", controllers.CashDeposit)
    router.GET("/transactions/history", controllers.TransactionHistory)

    router.Run()
}
