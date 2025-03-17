package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditRecord struct {
	UserUuid uuid.UUID `json:"userUuid"`
	CreateDt time.Time `json:"createDate"`
	UpdateDt time.Time `json:"updateDate"`
}

type ExerciseFields struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
	Cues           string `json:"cues"`
	PrimaryMuscles string `json:"primaryMuscles" validate:"required"`
	Apparatus      string `json:"apparatus"`
}

type ExerciseRequest struct {
	ExerciseFields
	UserUuid uuid.UUID `json:"userUuid"`
}

type Exercise struct {
	Uuid uuid.UUID `json:"uuid"`
	ExerciseFields
	AuditRecord
}
