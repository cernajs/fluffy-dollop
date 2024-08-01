package main

import (
	"fmt"
	"log"
	"os"

	handlers "testauth/auth"
	login "testauth/login"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	authParams := handlers.GoogleAuthParams{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		CallbackURL:  os.Getenv("CLIENT_CALLBACK_URL"),
	}

	if authParams.ClientID == "" || authParams.ClientSecret == "" || authParams.CallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	handlers.SetupAuth(authParams)

	dbParams := login.DBParams{
		DBPassword:       os.Getenv("DB_PASSWORD"),
		ConnectionString: fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/test", os.Getenv("DB_PASSWORD")),
	}

	if dbParams.DBPassword == "" || dbParams.ConnectionString == "" {
		log.Fatal("Environment variables (DB_PASSWORD, CONNECTION_STRING) are required")
	}

	login.SetupDB(dbParams)

	// google oauth
	r.GET("/", handlers.Home)
	r.GET("/auth/:provider", handlers.SignInWithProvider)
	r.GET("/auth/:provider/callback", handlers.CallbackHandler)
	r.GET("/success", handlers.Success)

	// login/register
	r.POST("/login", login.Login)

	r.GET("/register", login.Register)
	r.POST("/register-user", login.RegisterUser)

	r.Run(":5000")
}
