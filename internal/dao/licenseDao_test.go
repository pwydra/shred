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

func validLicense() *model.License {
	return &model.License{
		LicenseFields: model.LicenseFields{
			LicenseShortName: "CC-BY-SA 3",
			LicenseFullName:  "Creative Commons Attribution Share Alike 3",
			LicenseUrl:       "https://creativecommons.org/licenses/by-sa/3.0/deed.en",
		},
	}
}

func licenseReq() *model.LicenseRequest {
	return &model.LicenseRequest{
		LicenseFields: model.LicenseFields{
			LicenseShortName: "CC-BY-SA 3",
			LicenseFullName:  "Creative Commons Attribution Share Alike 3",
			LicenseUrl:       "https://creativecommons.org/licenses/by-sa/3.0/deed.en",
		},
	}
}

func validLicense2() *model.License {
	return &model.License{
		LicenseFields: model.LicenseFields{
			LicenseShortName: "CC0",
			LicenseFullName:  "Creative Commons Public Domain 1.0",
			LicenseUrl:       "http://creativecommons.org/publicdomain/zero/1.0/",
		},
	}
}

func invalidLicense() *model.License {
	return &model.License{
		LicenseFields: model.LicenseFields{
			LicenseShortName: "INVALID",
			LicenseFullName:  "Invalid",
			LicenseUrl:       "Invalid license",
		},
	}
}

func TestGetLicenseByShortName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	expected := validLicense()
	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseShortName := expected.LicenseShortName
	mock.ExpectQuery("SELECT \\* FROM license WHERE license_short_name = \\$1").
		WithArgs(licenseShortName).
		WillReturnRows(sqlmock.NewRows([]string{"license_short_name", "license_full_name", "url"}).
			AddRow(licenseShortName, expected.LicenseFullName, expected.LicenseUrl))

	license, err := dao.GetLicenseByShortName(licenseShortName)
	assert.NoError(t, err)
	assert.NotNil(t, license)
	assert.Equal(t, expected.LicenseFullName, license.LicenseFullName)
}

func TestGetLicenseByShortName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseShortName := "INVALID"
	mock.ExpectQuery("SELECT \\* FROM license WHERE license_short_name = \\$1").
		WithArgs(licenseShortName).
		WillReturnError(sql.ErrNoRows)

	license, err := dao.GetLicenseByShortName(licenseShortName)
	assert.Error(t, err)
	assert.Nil(t, license)
	assert.Equal(t, "license with short name INVALID not found", err.Error())
}

func TestGetAllLicenses(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	mock.ExpectQuery("SELECT \\* FROM license").
		WillReturnRows(sqlmock.NewRows([]string{"license_short_name", "license_full_name", "url"}).
			AddRow(validLicense().LicenseShortName, validLicense().LicenseFullName, validLicense().LicenseUrl).
			AddRow(validLicense2().LicenseShortName, validLicense2().LicenseFullName, validLicense2().LicenseUrl))

	ctx := context.Background()
	licenses, err := dao.GetAllLicenses(ctx)
	assert.NoError(t, err)
	assert.Len(t, licenses, 2)
	assert.Equal(t, validLicense().LicenseFullName, licenses[0].LicenseFullName)
	assert.Equal(t, validLicense2().LicenseFullName, licenses[1].LicenseFullName)
}

func TestCreateLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseReq := licenseReq()

	mock.ExpectExec("INSERT INTO license \\( license_short_name, license_full_name, url \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(licenseReq.LicenseShortName, licenseReq.LicenseFullName, licenseReq.LicenseUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.CreateLicense(licenseReq)
	assert.NoError(t, err)
}

func TestCreateLicense_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseReq := licenseReq()

	mock.ExpectExec("INSERT INTO license \\( license_short_name, license_full_name, url \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(licenseReq.LicenseShortName, licenseReq.LicenseFullName, licenseReq.LicenseUrl).
		WillReturnError(errors.New("insertion error"))

	err = dao.CreateLicense(licenseReq)
	assert.Error(t, err)
	assert.Equal(t, "insertion error", err.Error())
}

func TestUpdateLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseReq := licenseReq()

	mock.ExpectExec("UPDATE license SET license_full_name = \\$1, url = \\$2 WHERE license_short_name = \\$3").
		WithArgs(licenseReq.LicenseFullName, licenseReq.LicenseUrl, licenseReq.LicenseShortName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.UpdateLicense(licenseReq)
	assert.NoError(t, err)
}

func TestUpdateLicense_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseReq := licenseReq()

	mock.ExpectExec("UPDATE").
		WithArgs(licenseReq.LicenseFullName, licenseReq.LicenseUrl, licenseReq.LicenseShortName).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.UpdateLicense(licenseReq)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateLicense_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseReq := licenseReq()

	mock.ExpectExec("UPDATE license SET license_full_name = \\$1, url = \\$2 WHERE license_short_name = \\$3").
		WithArgs(licenseReq.LicenseFullName, licenseReq.LicenseUrl, licenseReq.LicenseShortName).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.UpdateLicense(licenseReq)
	assert.Error(t, err)
	assert.Equal(t, "license with Short Name 'CC-BY-SA 3' not found", err.Error())
}

func TestDeleteLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseShortName := "CC-BY-SA 3"
	mock.ExpectExec("DELETE FROM license WHERE license_short_name = \\$1").
		WithArgs(licenseShortName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.DeleteLicense(licenseShortName)
	assert.NoError(t, err)
}

func TestDeleteLicense_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseShortName := "CC-BY-SA 3"
	mock.ExpectExec("DELETE FROM license WHERE license_short_name = \\$1").
		WithArgs(licenseShortName).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.DeleteLicense(licenseShortName)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestDeleteLicense_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewLicenseDAO(sqlx.NewDb(db, "postgres"))

	licenseShortName := "INVALID"
	mock.ExpectExec("DELETE FROM license WHERE license_short_name = \\$1").
		WithArgs(licenseShortName).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.DeleteLicense(licenseShortName)
	assert.Error(t, err)
	assert.Equal(t, "license with Short Name 'INVALID' not found", err.Error())
}
