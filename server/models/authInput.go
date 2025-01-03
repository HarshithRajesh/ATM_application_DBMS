package models

type AuthInput struct{
  Username  string  `json:"username"  binding:"required"`
  Password  string  `json:"password"  binding:"required"`
    FirstName   string `json:"first_name" binding:"required"`
    LastName    string `json:"last_name" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginInput struct{
     Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"` 
}

