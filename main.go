package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Forceattack012/todoapidemo/auth"
	"github.com/Forceattack012/todoapidemo/todo"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println("Please consider environment variables %s", err)
	}

	dsn := os.Getenv("db")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&todo.Todo{})

	handler := todo.NewTodoHanlder(db)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	bytes := []byte(os.Getenv("SIGN"))
	r.GET("/tokenz", auth.AccessToken(bytes))

	proteced := r.Group("", auth.Protect(bytes))

	proteced.POST("/todos", handler.NewTask)
	r.GET("/todos/:id", handler.GetTodo)
	r.GET("/todos", handler.GetTodoList)
	r.Run()
}
