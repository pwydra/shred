package dao

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/pwydra/shred/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dbx := sqlx.NewDb(db, "postgres")
	defer db.Close()

	dao := NewExerciseDao(dbx)

	exReq := &model.ExerciseRequest{
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
		CreatedBy: uuid.New(),
	}

	mock.ExpectQuery("INSERT INTO exercise").
		WithArgs(exReq.ExerciseName, exReq.Description, exReq.Instructions, exReq.Cues,
			exReq.VideoUrl, exReq.CategoryCode, exReq.LicenseShortName, exReq.LicenseAuthor,
			exReq.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"exercise_uuid", "created_at", "updated_at"}).AddRow(uuid.New(), time.Now(), time.Now()))

	ex, err := dao.Create(exReq)
	assert.NoError(t, err)
	assert.NotNil(t, ex)
	assert.Equal(t, exReq.ExerciseName, ex.ExerciseName)
}

func TestCreateExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dbx := sqlx.NewDb(db, "postgres")
	defer db.Close()

	dao := NewExerciseDao(dbx)

	exReq := &model.ExerciseRequest{
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
		CreatedBy: uuid.New(),
	}

	mock.ExpectQuery("INSERT INTO exercise").
		WithArgs(exReq.ExerciseName, exReq.Description, exReq.Instructions, exReq.Cues,
			exReq.VideoUrl, exReq.CategoryCode, exReq.LicenseShortName, exReq.LicenseAuthor,
			exReq.CreatedBy).
		WillReturnError(sqlmock.ErrCancelled)

	ex, err := dao.Create(exReq)
	assert.Error(t, err)
	assert.Nil(t, ex)
	assert.Equal(t, "canceling query due to user request", err.Error())
}
func TestReadExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()
	createdAt := time.Now()

	mock.ExpectQuery("SELECT.*FROM exercise WHERE exercise_uuid =.*").
		WithArgs(exUuid).
		WillReturnRows(sqlmock.NewRows([]string{
			"exercise_uuid", "exercise_name", "exercise_description", "instructions",
			"cues", "video_url", "category_code", "license_short_name",
			"license_author", "created_by", "created_at", "updated_at"}).
			AddRow(
				exUuid, "Squat", "Lower Body", "standwith feet shoulder widthapart",
				"Keep chest up and back flat", "https://squat.mp4", "Strength", "MIT",
				"John Doe", uuid.New(), createdAt, createdAt))

	ex, err := dao.Read(exUuid)
	assert.NoError(t, err)
	assert.NotNil(t, ex)
	assert.Equal(t, "Squat", ex.ExerciseName)
}

func TestReadExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()

	mock.ExpectQuery("SELECT.*FROM exercise WHERE exercise_uuid =.*").
		WithArgs(exUuid).
		WillReturnError(sqlmock.ErrCancelled)

	ex, err := dao.Read(exUuid)
	assert.Error(t, err)
	assert.Nil(t, ex)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()

	ex := &model.Exercise{
		ExerciseUuid: exUuid,
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
	}

	mock.ExpectExec("UPDATE exercise SET.*WHERE exercise_uuid =.*").
		WithArgs(ex.ExerciseName, ex.Description, ex.Instructions, ex.Cues,
			ex.VideoUrl, ex.CategoryCode, ex.LicenseShortName, ex.LicenseAuthor,
			ex.ExerciseUuid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Update(ex)
	assert.NoError(t, err)
}

func TestUpdateExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()

	ex := &model.Exercise{
		ExerciseUuid: exUuid,
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
	}

	mock.ExpectExec("UPDATE exercise SET.*WHERE exercise_uuid =.*").
		WithArgs(ex.ExerciseName, ex.Description, ex.Instructions, ex.Cues,
			ex.VideoUrl, ex.CategoryCode, ex.LicenseShortName, ex.LicenseAuthor,
			ex.ExerciseUuid).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.Update(ex)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestDeleteExercise(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()
	mock.ExpectExec("DELETE FROM exercise WHERE exercise_uuid =.*").
		WithArgs(exUuid).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Delete(exUuid)
	assert.NoError(t, err)
}

func TestDeleteExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewExerciseDao(sqlx.NewDb(db, "postgres"))

	exUuid := uuid.New()
	mock.ExpectExec("DELETE FROM exercise WHERE exercise_uuid =.*").
		WithArgs(exUuid).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.Delete(exUuid)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}
