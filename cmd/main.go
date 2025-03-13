package main

import (
	"database/sql"
	"log"

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

func main() {
	log.Println("Starting Shred API")
	var err error
	connStr := "user=postgres dbname=shred_db password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
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
