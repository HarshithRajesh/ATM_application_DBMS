package controllers

import (
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/models"
  "net/http"
  "os"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
  "golang.org/x/crypto/bcrypt"
  "fmt"
  "strings"
)
func CreateUser(c *gin.Context) {
    var authInput models.AuthInput

    // Bind JSON input to struct
    if err := c.ShouldBindJSON(&authInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if all required fields are provided
    if authInput.Username == "" || authInput.Email == "" || authInput.Password == "" ||
       authInput.FirstName == "" || authInput.LastName == "" || authInput.PhoneNumber == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required: username, email, password, first name, last name, phone number"})
        return
    }

    // Check if the username already exists
    var userFound models.User
    initializers.DB.Where("username = ?", authInput.Username).First(&userFound)
    if userFound.UserID != 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
        return
    }

    // Check if the email already exists
    initializers.DB.Where("email = ?", authInput.Email).First(&userFound)
    if userFound.UserID != 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
        return
    }

    // Hash the password
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create the user
    user := models.User{
        Username:   authInput.Username,
        Password:   string(passwordHash),
        FirstName:  authInput.FirstName,
        LastName:   authInput.LastName,
        Email:      authInput.Email,
        PhoneNumber: authInput.PhoneNumber,
    }

    // Save the user to the database
    if err := initializers.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Return the created user (excluding password)
    user.Password = "" // Don't return the password in the response
    c.JSON(http.StatusOK, gin.H{
        "data": user,
    })
}


func Login(c *gin.Context) {
    var authInput models.LoginInput

    // Bind JSON input to struct for login
    if err := c.ShouldBindJSON(&authInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var userFound models.User
    if authInput.Username != "" {
        initializers.DB.Where("username = ?", authInput.Username).First(&userFound)
    }
    if userFound.UserID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
        return
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
        return
    }

    // Generate JWT token
    generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":  userFound.UserID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    })
    token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
    })
}


func GetUserProfile(c *gin.Context){
  user, _ := c.Get("currentUser")
  u := user.(User
  c.JSON(200,gin.H{
    "user":user,
  })
}

func CheckAuth(c *gin.Context){

  authHeader := c.GetHeader("Authorization")

  if authHeader == ""{
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Authorization header is missing"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  authToken := strings.Split(authHeader," ")
  if len(authToken) != 2 || authToken[0] != "Bearer" {
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid Token format"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  tokenString := authToken[1]
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
      return nil,fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
    }
    return []byte(os.Getenv("SECRET")),nil
  })
  
  if err != nil || !token.Valid {
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid or expired token"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }
  claims, ok := token.Claims.(jwt.MapClaims)
  if !ok {
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid token"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  if float64(time.Now().Unix()) > claims["exp"].(float64){
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Token expired"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }
  var user models.User
  initializers.DB.Where("ID=?",claims["id"]).Find(&user)

  if user.UserID == 0 {
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }
  c.Set("currentUser",user)
  c.Next()
}

