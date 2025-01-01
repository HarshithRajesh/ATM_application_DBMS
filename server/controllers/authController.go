package controller

import (
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/models"
  "net/http"
  "os"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
  "golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context){
  var authInput models.authInput

  if err := c.ShouldBindJSON(&authInput); err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
    return
  }

  var userFound models.User 
  initializers.DB.Where("username=?",authInput.Username).Find(&userFound)

  if userFound.ID != 0{
    c.JSON(http.StatusBadRequest,gin.H{"error":"Username already exists"})
    return
  }
  passwordHash,err := bcrypt.GenerateFromPassword([]byte(authInput.Password),bcrypt.DefaultCost)
  if err != nil{
    c.JSON(http.StatusBadRequest,gin.H{"error":err.Error})
    return
  }
  user := models.User{
    Username: authInput.Username,
    Password: string(passwordHash),
  }

  initializers.DB.Create(&user)
  c.JSON(http.StatusOK,gin.H{"data":user})
  
}

func Login(c *gin.Context){
  var authInput models.authInput

  if err := c.ShouldBindJSON(&authInput); err != nil{
    c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
    return
  }

  var userFound models.User
  initializers.DB.Where("username=?",authInput.Username).Find(&userFound)

  if userFound.ID == 0{
    c.JSON(http.StatusBadRequest,gin.H{"error":"user not found"})
    return
  }

  if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password),[]byte(authInput.Password)); err != nil{
    c.JSON(https.StatusBadRequest,gin.H{"error":"invalid password"})
    return
  }

  generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "id": userFound.UserId,
    "exp": time.Now().Add(time.Hour*24).Unix(),
  })
  token, err :=   generateToken.SignedString([]byte(os.Getenv("SECRET")))

  if err != nil {
    c.JSON(http.StatusBadRequest,gin.H{"error":"Failed to generate token"})
  }

  c.JSON(200,gin.H{
    "token":token,
  })
}

func GetUserProfile(c *gin.Context){
  user, _ := c.Get("currentUser")
  c.JSON(200,gin.H{
    "user":user,
  })
}

func CheckAuth(c *gin.Context){

  authHeader := c.GetHeader("Authorization")

  if authHeader == ""{
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Authorization header is missing"})
    c.AbortWithStatus(http.StausUnauthorized)
    return
  }

  authToken := strings.Split(authHeader," ")
  if len(authToken) != 2 || authToken[0] != "Bearer" {
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid Token format"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  tokenString := authToken[1]
  token, err := jwt.Parse(tokenString, func(token, *jwt.Token) (interface{},error){
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
  claims, ok := token.Claims(jwt.MapClaims)
  if !ok {
    c.JSON(http.StausUnauthorized,gin.H{"error":"Invalid token"})
    c.Abort()
    return
  }

  if float64(time.Now().Unix()) > claims["exp"].(float64){
    c.JSON(http.StatusUnauthorized,gin.H{"error":"Token expired"})
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }
  var user models.User
  initializers.DB.Where("ID=?",claims["id"]).Find(&user)

  if user.ID == 0 {
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }
  c.Set("currentUser",user)
  c.Next()
}

