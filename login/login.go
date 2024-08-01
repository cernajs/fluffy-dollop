package login

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type DBParams struct {
	DBPassword       string
	ConnectionString string
}

var dbParams DBParams

func SetupDB(params DBParams) {
	dbParams = params
}

func Login(c *gin.Context) {
	db, err := sql.Open("mysql", dbParams.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	var name string
	err = db.QueryRow("SELECT name FROM users WHERE email = ? AND password = ?", email, password).Scan(&name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//c.JSON(http.StatusOK, gin.H{"name": name})

	c.Redirect(http.StatusFound, "/success")

	defer db.Close()
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
	db, err := sql.Open("mysql", dbParams.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Redirect(http.StatusFound, "/success")

	defer db.Close()
}
