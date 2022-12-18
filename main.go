package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Forceattack012/todoapidemo/router"
	"github.com/Forceattack012/todoapidemo/store"
	"github.com/Forceattack012/todoapidemo/todo"
)

var (
	buildCommit = "dev"
	buildTime   = time.Now().String()
)

func main() {
	_, err := os.Create("/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("Please consider environment variables %s", err.Error())
	}

	dsn := os.Getenv("db")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&todo.Todo{})

	gormStore := store.NewGormStore(db)
	handler := todo.NewTodoHanlder(gormStore)
	r := router.NewFiberRouter()

	r.GET("/todos/:id", handler.GetTodo)
	r.GET("/todos", handler.GetTodoList)
	r.POST("/todos", handler.NewTask)
	r.PATCH("/todos/:id", handler.UpdateTask)
	r.DELETE("/todos/:id", handler.RemoveTask)

	r.Listen(":" + os.Getenv("PORT"))

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// r.GET("/healthz", func(c *gin.Context) {
	// 	c.Status(200)
	// })

	// // r.GET("limitx", limiterHandler)

	// r.GET("/x", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		buildCommit: buildCommit,
	// 		buildTime:   buildTime,
	// 	})
	// })

	//bytes := []byte(os.Getenv("SIGN"))
	// r.GET("/tokenz", auth.AccessToken(bytes))

	// proteced := r.Group("", auth.Protect(bytes))

	// proteced.GET("/todos/:id", router.NewGinContext(handler.GetTodo))
	// proteced.GET("/todos", router.NewGinContext(handler.GetTodoList))
	// proteced.POST("/todos", router.NewGinContext(handler.NewTask))
	// proteced.PATCH("/todos/:id", router.NewGinContext(handler.UpdateTask))
	// proteced.DELETE("/todos/:id", router.NewGinContext(handler.RemoveTask))

	// create notify context for recived signal SIGINT or SIGTERM
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	// s := &http.Server{
	// 	Addr:           ":" + os.Getenv("PORT"),
	// 	Handler:        r,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// go func() {
	// 	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("listen : %s\n", err)
	// 	}
	// }()

	// //when signal sigint or sigterm input it'll stop serivce
	// <-ctx.Done()
	// stop()
	// fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	// timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := s.Shutdown(timeoutContext); err != nil {
	// 	fmt.Println(err)
	// }
}

// var limiter = rate.NewLimiter(5, 5)

// func limiterHandler(c *gin.Context) {
// 	if !limiter.Allow() {
// 		c.AbortWithStatus(http.StatusTooManyRequests)
// 	}

// 	c.JSON(200, gin.H{
// 		"message": "pong",
// 	})
// }
