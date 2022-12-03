package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Forceattack012/todoapidemo/auth"
	"github.com/Forceattack012/todoapidemo/todo"
)

var (
	buildCommit = "dev"
	buildTime   = time.Now().String()
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

	r.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			buildCommit: buildCommit,
			buildTime:   buildTime,
		})
	})

	bytes := []byte(os.Getenv("SIGN"))
	r.GET("/tokenz", auth.AccessToken(bytes))

	proteced := r.Group("", auth.Protect(bytes))

	proteced.POST("/todos", handler.NewTask)
	r.GET("/todos/:id", handler.GetTodo)
	r.GET("/todos", handler.GetTodoList)

	// create notify context for recived signal SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen : %s\n", err)
		}
	}()

	//when signal sigint or sigterm input it'll stop serivce
	<-ctx.Done()
	stop()
	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutContext); err != nil {
		fmt.Println(err)
	}
}
