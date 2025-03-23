package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
		WillReturnRows(sqlmock.NewRows([]string{"muscle_code", "muscle_name", "muscle_description", "muscle_group"}).
			AddRow(musCode, "Latissimus", "Muscle of the back", "Back"))

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
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode: "LAT",
			MuscleName: "Latissimus",
			MuscleDesc: "Muscle of the back",
		},
	}

	mock.ExpectExec("INSERT INTO muscle_type \\( muscle_code, muscle_name, muscle_description, muscle_group \\) VALUES \\( \\$1, \\$2, \\$3, \\$4 \\)").
		WithArgs(catReq.MuscleCode, catReq.MuscleName, catReq.MuscleDesc).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.CreateMuscle(catReq)
	assert.NoError(t, err)
}

func TestCreateMuscle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode: "LAT",
			MuscleName: "Latissimus",
			MuscleDesc: "Muscle of the back",
		},
	}

	mock.ExpectExec("INSERT INTO muscle_type \\( muscle_code, muscle_name, muscle_description, muscle_group \\) VALUES \\( \\$1, \\$2, \\$3, \\$4 \\)").
		WithArgs(catReq.MuscleCode, catReq.MuscleName, catReq.MuscleDesc).
		WillReturnError(errors.New("insertion error"))

	err = dao.CreateMuscle(catReq)
	assert.Error(t, err)
	assert.Equal(t, "insertion error", err.Error())
}

func TestUpdateMuscle(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "LAT",
			MuscleName:  "Latissimus",
			MuscleDesc:  "Muscle of the back",
			MuscleGroup: "Back",
		},
	}

	mock.ExpectExec("UPDATE muscle_type SET muscle_name = \\$1, muscle_description = \\$2, muscle_group = \\$3 WHERE muscle_code = \\$4").
		WithArgs(catReq.MuscleName, catReq.MuscleDesc, catReq.MuscleGroup, catReq.MuscleCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.UpdateMuscle(catReq)
	assert.NoError(t, err)
}

func TestUpdateMuscle_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "LAT",
			MuscleName:  "Latissimus",
			MuscleDesc:  "Muscle of the back",
			MuscleGroup: "Back",
		},
	}

	mock.ExpectExec("UPDATE").
		WithArgs(catReq.MuscleName, catReq.MuscleDesc, catReq.MuscleGroup, catReq.MuscleCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.UpdateMuscle(catReq)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateMuscle_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewMuscleDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.MuscleRequest{
		MuscleFields: model.MuscleFields{
			MuscleCode:  "INVALID",
			MuscleName:  "Invalid",
			MuscleDesc:  "Invalid description",
			MuscleGroup: "Invalid group",
		},
	}

	mock.ExpectExec("UPDATE muscle_type SET muscle_name = \\$1, muscle_description = \\$2, muscle_group = \\$3 WHERE muscle_code = \\$4").
		WithArgs(catReq.MuscleName, catReq.MuscleDesc, catReq.MuscleGroup, catReq.MuscleCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.UpdateMuscle(catReq)
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
