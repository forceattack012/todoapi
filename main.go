package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Forceattack012/todoapidemo/todo"
)

func main() {
	dsn := "root:123456789@tcp(127.0.0.1:3306)/dogapp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
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

	r.POST("/todos", handler.NewTask)
	r.GET("/todos/:id", handler.GetTodo)
	r.GET("/todos", handler.GetTodoList)
	r.Run()
}
