package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pwydra/shred/internal/dao"
	"github.com/pwydra/shred/internal/model"
)

type Handler struct {
	dao dao.ExerciseDaoInterface
}

func NewHandler(dao dao.ExerciseDaoInterface) *Handler {
	return &Handler{dao: dao}
}

func (h Handler) CreateExercise(ctx *gin.Context) {
	var exReq model.ExerciseRequest
	if err := ctx.ShouldBindJSON(&exReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ex, err := h.dao.Create(&exReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ex)
}

func (h Handler) UpdateExercise(ctx *gin.Context) {
	var exReq model.Exercise
	if err := ctx.ShouldBindJSON(&exReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exReq.ExerciseUuid != uuid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UUID in path does not match UUID in request body"})
		return
	}

	err = h.dao.Update(&exReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, exReq)
}

func (h Handler) DeleteExercise(ctx *gin.Context) {
	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.dao.Delete(uuid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (h Handler) GetExercise(ctx *gin.Context) {
	uuid, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ex, err := h.dao.Read(uuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, ex)
}
