package dao

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/pwydra/shred/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateExercise(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    dao := NewExerciseDao(db)

    exReq := &model.ExerciseRequest{
        ExerciseFields: model.ExerciseFields{
            Name:           "Squat",
            Description:    "Lower Body",
            Cues:           "Keep chest up and back flat",
            PrimaryMuscles: "Quadriceps",
            Apparatus:      "Barbell",
        },
		UserUuid: uuid.New(),
    }

    mock.ExpectQuery("INSERT INTO exercise").
        WithArgs(exReq.Name, exReq.Description, exReq.Cues, exReq.PrimaryMuscles, exReq.Apparatus, exReq.UserUuid).
        WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(uuid.New()))

    ex, err := dao.Create(exReq)
    assert.NoError(t, err)
    assert.NotNil(t, ex)
    assert.Equal(t, exReq.Name, ex.Name)
}

func TestReadExercise(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    dao := NewExerciseDao(db)

    exUuid := uuid.New()
	createdAt := time.Now()

    mock.ExpectQuery("SELECT.*FROM exercise WHERE uuid =.*").
        WithArgs(exUuid).
        WillReturnRows(sqlmock.NewRows([]string{"uuid", "exercise_name", "description", "cues", "primary_muscles", "apparatus", "created_at", "user_uuid"}).
            AddRow(exUuid, "Squat", "Lower Body", "Keep chest up and back flat", "Quadriceps", "Barbell", createdAt, uuid.New()))

    ex, err := dao.Read(exUuid)
    assert.NoError(t, err)
    assert.NotNil(t, ex)
    assert.Equal(t, "Squat", ex.Name)
}

func TestUpdateExercise(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    dao := NewExerciseDao(db)

	exUuid := uuid.New()

    ex := &model.Exercise{
        Uuid:           exUuid,
        ExerciseFields: model.ExerciseFields{
			Name: "Squat", 
			Description: "Lower Body", 
			Cues: "Keep chest up and back flat", 
			PrimaryMuscles: "Quadriceps", 
			Apparatus: "Barbell",
		},
    }

    mock.ExpectExec("UPDATE exercise SET.*WHERE uuid =.*").
        WithArgs(ex.Name, ex.Description, ex.Cues, ex.PrimaryMuscles, ex.Apparatus, ex.Uuid).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = dao.Update(ex)
    assert.NoError(t, err)
}

func TestDeleteExercise(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    dao := NewExerciseDao(db)

    exUuid := uuid.New()
    mock.ExpectExec("DELETE FROM exercise WHERE uuid =.*").
        WithArgs(exUuid).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = dao.Delete(exUuid)
    assert.NoError(t, err)
}