package middleware

import(
  "github.com/HarshithRajesh/zapster/initializers"
  "github.com/HarshithRajesh/zapster/initializers"
  "net/http"
  "os"
  "time"
  "strings"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
)

func checkAuth(c *gin.Context){

  authHeader := c.GetHeader("Authorization")

  if authHeader == ""{
    c.JSON(http.StatusUnauthorized,gin.H{
      "error":"Authorization header is missing"
    })
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  authToken := strings.Split(authHeader," ")
  if len(authToken)!=2 || authToken[0] != 'Bearer'{
    c.JSON(http.StatusUnauthorized,gin.H{
      "error":"Invalid token format"
    })
    c.AbortWithStatus(http.StatusUnauthorized)
    return
  }

  tokenString := authToken[1]
  token , err := jwt.Parse(tokenString,func(token *jwt.Token) (interface{},error){
    if _,ok :=token.Method.(*jwt.SigningMethodHMAC); !ok{
      return nil,fmt.Errorf("Unexpected signing method: %v",token.Header["alg"])
    }
    return []byte(os.Getenv("SECRET")),nil

})
if err
}
