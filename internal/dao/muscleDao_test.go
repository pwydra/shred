package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pwydra/shred/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetMuscleByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musCode := "LATISSIMUS"
	mock.ExpectQuery("SELECT \\* FROM muscle_type WHERE muscle_code = \\$1").
		WithArgs(musCode).
		WillReturnRows(sqlmock.NewRows([]string{"muscle_code", "muscle_name", "muscle_description", "muscle_group", "created_by", "created_at", "updated_at"}).
			AddRow(musCode, "Latissimus", "Muscle of the back", "Back", uuid.New(), time.Now(), time.Now()))

	muscle, err := dao.GetMuscleByCode(musCode)
	assert.NoError(t, err)
	assert.NotNil(t, muscle)
	assert.Equal(t, "Latissimus", muscle.MuscleName)
}

func TestGetMuscleByCode_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catCode := "INVALID"
	mock.ExpectQuery("SELECT \\* FROM muscle_type WHERE muscle_code = \\$1").
		WithArgs(catCode).
		WillReturnError(sql.ErrNoRows)

	muscle, err := dao.GetMuscleByCode(catCode)
	assert.Error(t, err)
	assert.Nil(t, muscle)
	assert.Equal(t, "muscle with code INVALID not found", err.Error())
}

func TestGetAllMuscles(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	mock.ExpectQuery("SELECT \\* FROM muscle_type").
		WillReturnRows(sqlmock.NewRows([]string{"muscle_code", "muscle_name", "muscle_description", "muscle_group"}).
			AddRow("LAT", "Latissimus", "Muscle of the back", "Back").
			AddRow("BICEP", "Bicep", "Muscle of the upper arm", "Arm"))

	ctx := context.Background()
	categories, err := dao.GetAllMuscles(ctx)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Latissimus", categories[0].MuscleName)
	assert.Equal(t, "Bicep", categories[1].MuscleName)
}

func TestCreateMuscle(t *testing.T) {
	timeNow := time.Now()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode: "LAT",
			MuscleName: "Latissimus",
			MuscleDesc: "Muscle of the back",
			MuscleGroup: "Back",
		},
		CreatedBy: uuid.New(),
	}

	mock.ExpectQuery("INSERT INTO muscle_type \\( muscle_code, muscle_name, muscle_description, muscle_group, created_by \\) VALUES \\( \\$1, \\$2, \\$3, \\$4, \\$5 \\).*").
		WithArgs(musReq.MuscleCode, musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, musReq.CreatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(timeNow, timeNow))

	mus, err := dao.CreateMuscle(musReq)
	assert.NoError(t, err)
	assert.Equal(t, mus.MuscleCode, musReq.MuscleCode)
	assert.Equal(t, mus.MuscleName, musReq.MuscleName)
	assert.Equal(t, mus.MuscleDesc, musReq.MuscleDesc)
	assert.Equal(t, mus.CreatedAt, timeNow)
	assert.Equal(t, mus.UpdatedAt, timeNow)
}

func TestCreateMuscle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode: "LAT",
			MuscleName: "Latissimus",
			MuscleDesc: "Muscle of the back",
		},
	}

	mock.ExpectQuery("INSERT INTO muscle_type \\( muscle_code, muscle_name, muscle_description, muscle_group, created_by \\) VALUES \\( \\$1, \\$2, \\$3, \\$4, \\$5 \\).*").
		WithArgs(musReq.MuscleCode, musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, musReq.CreatedBy).
		WillReturnError(errors.New("insertion error"))

	_, err = dao.CreateMuscle(musReq)
	assert.Error(t, err)
	assert.Equal(t, "insertion error", err.Error())
}

func TestUpdateMuscle(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "LAT",
			MuscleName:  "Latissimus",
			MuscleDesc:  "Muscle of the back",
			MuscleGroup: "Back",
		},
	}

	mock.ExpectExec("UPDATE muscle_type SET muscle_name = \\$1, muscle_description = \\$2, muscle_group = \\$3 WHERE muscle_code = \\$4").
		WithArgs(musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, musReq.MuscleCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.UpdateMuscle(musReq)
	assert.NoError(t, err)
}

func TestUpdateMuscle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "LAT",
			MuscleName:  "Latissimus",
			MuscleDesc:  "Muscle of the back",
			MuscleGroup: "Back",
		},
	}

	mock.ExpectExec("UPDATE").
		WithArgs(musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, musReq.MuscleCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.UpdateMuscle(musReq)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateMuscle_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	musReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "INVALID",
			MuscleName:  "Invalid",
			MuscleDesc:  "Invalid description",
			MuscleGroup: "Invalid group",
		},
	}

	mock.ExpectExec("UPDATE muscle_type SET muscle_name = \\$1, muscle_description = \\$2, muscle_group = \\$3 WHERE muscle_code = \\$4").
		WithArgs(musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, musReq.MuscleCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.UpdateMuscle(musReq)
	assert.Error(t, err)
	assert.Equal(t, "muscle with Code INVALID not found", err.Error())
}

func TestDeleteMuscle(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catCode := "STRENGTH"
	mock.ExpectExec("DELETE FROM muscle_type WHERE muscle_code = \\$1").
		WithArgs(catCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.DeleteMuscle(catCode)
	assert.NoError(t, err)
}

func TestDeleteMuscle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catCode := "STRENGTH"
	mock.ExpectExec("DELETE FROM muscle_type WHERE muscle_code = \\$1").
		WithArgs(catCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.DeleteMuscle(catCode)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestDeleteMuscle_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catCode := "INVALID"
	mock.ExpectExec("DELETE FROM muscle_type WHERE muscle_code = \\$1").
		WithArgs(catCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.DeleteMuscle(catCode)
	assert.Error(t, err)
	assert.Equal(t, "muscle with code INVALID not found", err.Error())
}
