package main

import (
	"fmt"
	"log"
	"os"
	"peopledb/controllers"
	"peopledb/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	gin.DefaultWriter = infoLog.Writer()
	err := godotenv.Load(".env")
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"))
	err = models.ConnectDB(dsn)
	if err != nil {
		errorLog.Fatal()
		return
	}

	r.POST("/user", controllers.AddUser)
	r.GET("/user", controllers.GetUser)
	r.GET("/user/all", controllers.GetAllUsers)
	r.PATCH("/user", controllers.UpdateUser)
	r.DELETE("/user", controllers.DeleteUser)
	r.Run(os.Getenv("SERVER_ADDR"))
}
