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

func TestGetApparatusByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appCode := "BENCH"
	mock.ExpectQuery("SELECT \\* FROM apparatus_type WHERE apparatus_code = \\$1").
		WithArgs(appCode).
		WillReturnRows(sqlmock.NewRows([]string{"apparatus_code", "apparatus_name", "apparatus_description"}).
			AddRow(appCode, "Bench", "Bench you can lie on and that is optionally adjustable"))

	apparatus, err := dao.GetApparatusByCode(appCode)
	assert.NoError(t, err)
	assert.NotNil(t, apparatus)
	assert.Equal(t, "Bench", apparatus.ApparatusName)
}

func TestGetApparatusByCode_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appCode := "INVALID"
	mock.ExpectQuery("SELECT \\* FROM apparatus_type WHERE apparatus_code = \\$1").
		WithArgs(appCode).
		WillReturnError(sql.ErrNoRows)

	apparatus, err := dao.GetApparatusByCode(appCode)
	assert.Error(t, err)
	assert.Nil(t, apparatus)
	assert.Equal(t, "apparatus with code INVALID not found", err.Error())
}

func TestGetAllApparatuses(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	mock.ExpectQuery("SELECT \\* FROM apparatus_type").
		WillReturnRows(sqlmock.NewRows([]string{"apparatus_code", "apparatus_name", "apparatus_description"}).
			AddRow("BENCH", "Bench", "Bench you can lie on and that is optionally adjustable").
			AddRow("BAND", "Band", "Elastic band for resistance training"))

	ctx := context.Background()
	apparatuses, err := dao.GetAllApparatuses(ctx)
	assert.NoError(t, err)
	assert.Len(t, apparatuses, 2)
	assert.Equal(t, "Bench", apparatuses[0].ApparatusName)
	assert.Equal(t, "Band", apparatuses[1].ApparatusName)
}

func TestCreateApparatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appReq := &model.ApparatusRequest{
		ApparatusFields: model.ApparatusFields{
		ApparatusCode: "BENCH",
		ApparatusName: "Bench",
		ApparatusDesc: "Bench you can lie on and that is optionally adjustable",
		},
	}

	mock.ExpectExec("INSERT INTO apparatus_type \\( apparatus_code, apparatus_name, apparatus_description \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(appReq.ApparatusCode, appReq.ApparatusName, appReq.ApparatusDesc).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.CreateApparatus(appReq)
	assert.NoError(t, err)
}

func TestCreateApparatus_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appReq := &model.ApparatusRequest{
		ApparatusFields: model.ApparatusFields{
			ApparatusCode: "BENCH",
			ApparatusName: "Bench",
			ApparatusDesc: "Bench you can lie on and that is optionally adjustable",
			},
		}

	mock.ExpectExec("INSERT INTO apparatus_type \\( apparatus_code, apparatus_name, apparatus_description \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(appReq.ApparatusCode, appReq.ApparatusName, appReq.ApparatusDesc).
		WillReturnError(errors.New("insertion error"))

	err = dao.CreateApparatus(appReq)
	assert.Error(t, err)
	assert.Equal(t, "insertion error", err.Error())
}

func TestUpdateApparatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appReq := &model.ApparatusRequest{
		ApparatusFields: model.ApparatusFields{
			ApparatusCode: "BENCH",
			ApparatusName: "Bench",
			ApparatusDesc: "Bench you can lie on and that is optionally adjustable",
		},
	}

	mock.ExpectExec("UPDATE apparatus_type SET apparatus_name = \\$1, apparatus_description = \\$2 WHERE apparatus_code = \\$3").
		WithArgs(appReq.ApparatusName, appReq.ApparatusDesc, appReq.ApparatusCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.UpdateApparatus(appReq)
	assert.NoError(t, err)
}

func TestUpdateApparatus_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appReq := &model.ApparatusRequest{
		ApparatusFields: model.ApparatusFields{
			ApparatusCode: "BENCH",
			ApparatusName: "Bench",
			ApparatusDesc: "Bench you can lie on and that is optionally adjustableBench you can lie on and that is optionally adjustable",
		},
	}

	mock.ExpectExec("UPDATE").
	WithArgs(appReq.ApparatusName, appReq.ApparatusDesc, appReq.ApparatusCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.UpdateApparatus(appReq)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateApparatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appReq := &model.ApparatusRequest{
		ApparatusFields: model.ApparatusFields{
			ApparatusCode: "INVALID",
			ApparatusName: "Invalid",
			ApparatusDesc: "Invalid description",
		},
	}

	mock.ExpectExec("UPDATE apparatus_type SET apparatus_name = \\$1, apparatus_description = \\$2 WHERE apparatus_code = \\$3").
		WithArgs(appReq.ApparatusName, appReq.ApparatusDesc, appReq.ApparatusCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.UpdateApparatus(appReq)
	assert.Error(t, err)
	assert.Equal(t, "apparatus with Code INVALID not found", err.Error())
}

func TestDeleteApparatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appCode := "BENCH"
	mock.ExpectExec("DELETE FROM apparatus_type WHERE apparatus_code = \\$1").
		WithArgs(appCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.DeleteApparatus(appCode)
	assert.NoError(t, err)
}

func TestDeleteApparatus_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appCode := "BENCH"
	mock.ExpectExec("DELETE FROM apparatus_type WHERE apparatus_code = \\$1").
		WithArgs(appCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.DeleteApparatus(appCode)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}


func TestDeleteApparatus_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewApparatusDAO(sqlx.NewDb(db, "postgres"))

	appCode := "INVALID"
	mock.ExpectExec("DELETE FROM apparatus_type WHERE apparatus_code = \\$1").
		WithArgs(appCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.DeleteApparatus(appCode)
	assert.Error(t, err)
	assert.Equal(t, "apparatus with code INVALID not found", err.Error())
}