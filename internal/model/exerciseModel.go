package model

import (
	"time"

	"github.com/google/uuid"
)

/*
 * The AuditRecord struct is used to track the creation and last update time of a record.
 */
type AuditRecord struct {
	CreatedBy uuid.UUID `json:"createdBy" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type MuscleFields struct {
	MuscleCode  string `json:"muscleCode" db:"muscle_code"`
	MuscleName  string `json:"muscleName" db:"muscle_name"`
	MuscleDesc  string `json:"muscleDesc" db:"muscle_description"`
	MuscleGroup string `json:"muscleGroup" db:"muscle_group"`
}

type MuscleRequest struct {
	MuscleFields
}

type Muscle struct {
	MuscleFields
	AuditRecord
}

type CategoryFields struct {
	CategoryCode string `json:"categoryCode" db:"category_code"`
	CategoryName string `json:"categoryName" db:"category_name"`
	CategoryDesc string `json:"categoryDesc" db:"category_description"`
}

type CategoryRequest struct {
	CategoryFields
}

type Category struct {
	CategoryFields
	AuditRecord
}

type LicenseFields struct {
	LicenseShortName string `json:"licenseShortName"`
	LicenseFullName  string `json:"licenseFullName"`
	LicenseUrl       string `json:"licenseUrl"`
}

type LicenseRequest struct {
	LicenseFields
}

type License struct {
	LicenseFields
	AuditRecord
}

type ExerciseFields struct {
	ExerciseName     string `json:"exerciseName" validate:"required" db:"exercise_name"`
	Description      string `json:"description" db:"exercise_description"`
	Instructions     string `json:"instructions" db:"instructions"`
	Cues             string `json:"cues" db:"cues"`
	VideoUrl         string `json:"videoUrl" db:"video_url"`
	CategoryCode     string `json:"category" db:"category_code"`
	LicenseShortName string `json:"licenceShortName" db:"license_short_name"`
	LicenseAuthor    string `json:"licenceAuthor" db:"license_author"`
}

type ExerciseRequest struct {
	ExerciseFields
	CreatedBy uuid.UUID `json:"createdBy" db:"created_by"`
}

type Exercise struct {
	ExerciseUuid uuid.UUID `json:"exerciseUuid" db:"exercise_uuid"`
	ExerciseFields
	AuditRecord
}
