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

	"github.com/chakhrits/todo-api/auth"
	"github.com/chakhrits/todo-api/todo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	//get environment from env
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variable: %s", err)
	}
	//db connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//Automigrate
	db.AutoMigrate(&todo.Todo{})

	r := gin.Default() //Create default router
	//handler with path
	r.GET("/ping", pingpongHandler)
	handler := todo.NewTodoHandler(db)

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))
	//สร้างgroupเพื่อแยกว่าpathไหนprotectได้บ้าง
	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	protected.POST("/todos", handler.NewTask)

	//Graceful shutdown
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
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down graceful, press CTRL+C again to force")

	// r.Run() //set port default is :8080
}

func pingpongHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
