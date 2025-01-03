package main

import (
  "log"
  "github.com/HarshithRajesh/zapster/initializers"

  "github.com/HarshithRajesh/zapster/models"
)

func init () {
  initializers.LoadEnvs()
  initializers.ConnectDB()
}

func main(){
 log.Println("Starting migration...")
    err := initializers.DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    } else {
        log.Println("Migration completed successfully")
    }
}
