package main

import (
  "github.com/HarshithRajesh/zapster/initializers"

  "github.com/HarshithRajesh/zapster/models"
)

func init () {
  initializers.LoadEnvs()
  initializers.ConnectDB()
}

func main(){
  initializers.DB.AutoMigrate(&models.User{})
}
