package login

import (
	"context"
	"log"
	"net/http"
	db "testauth/db"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	//"golang.org/x/crypto/bcrypt"
)

var dbConnection db.IDBconnection

func SetDBConnection(db db.IDBconnection) {
	dbConnection = db
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	//hash, err := bcrypt.GenerateFromPassword([]byte(password),10)

	ctx := context.Background()
	//name, err := dbConnection.FindUserByEmailAndPassword(ctx, email, password)
	_, err := dbConnection.FindUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret")) //os.Getenv("SECRET"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	c.Redirect(http.StatusFound, "/success")
}

func Register(c *gin.Context) {
	log.Default().Println("Registering user")
	res, err := template.ParseFiles("templates/register.html")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = res.Execute(c.Writer, gin.H{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func RegisterUser(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	ctx := context.Background()
	err := dbConnection.CreateUser(ctx, name, email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/success")
}

func RequiresAuth(c *gin.Context) {

	cookie, err := c.Cookie("token")
	log.Default().Println("Cookie: ", cookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	log.Default().Println("Claims: ", claims)
	log.Default().Println("OK: ", ok)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// You can set claims to context if needed
	c.Set("email", claims["email"])
	c.Next()
}
