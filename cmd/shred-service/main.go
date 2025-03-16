package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/pwydra/shred/internal/dao"
	"github.com/pwydra/shred/internal/handlers"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}

func getConnectionString() string {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Database environment variables are not set")
	}

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, port, dbUser, dbPassword, dbName)
}

func main() {
	log.Println("Starting Shred API")

	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	exerciseDao := dao.NewExerciseDao(db)
	handler := handlers.NewHandler(exerciseDao)

	r := NewRouter()

	//	r.Engine.GET("/exercises/query", handler.queryExercise)
	//	r.Engine.GET("/exercises", handler.GetExercises)
	r.Engine.GET("/exercises/:uuid", handler.GetExercise)
	r.Engine.POST("/exercises", handler.CreateExercise)
	r.Engine.PUT("/exercises/:uuid", handler.UpdateExercise)
	r.Engine.DELETE("/exercises/:uuid", handler.DeleteExercise)

	if err := r.Engine.Run(":8088"); err != nil {
		panic(err)
	}
}
