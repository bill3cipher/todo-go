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
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	//Liveness probe check api
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	//get environment from env
	err = godotenv.Load("local.env")
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
	//===============Router===============
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("limit", limitedHandler)
	r.GET("/x", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	//handler with path
	r.GET("/ping", pingpongHandler)
	handler := todo.NewTodoHandler(db)

	r.GET("/token", auth.AccessToken(os.Getenv("SIGN")))
	//สร้างgroupเพื่อแยกว่าpathไหนprotectได้บ้าง
	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	protected.POST("/todos", handler.NewTask)
	//===============End Router===============

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

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
