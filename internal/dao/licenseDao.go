package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pwydra/shred/internal/model"
)

// LicenseDAO provides access to the licenses in the database.
type LicenseDAO struct {
	db *sqlx.DB
}

// NewLicenseDAO creates a new instance of LicenseDAO.
func NewLicenseDAO(db *sqlx.DB) *LicenseDAO {
	return &LicenseDAO{db: db}
}

const createLicenseDML string = `
	INSERT INTO license (
		license_short_name, license_full_name, url
	) VALUES (
		$1, $2, $3
	)`

// GetLicenseByShortName retrieves a license by its short name.
const getLicenseByShortNameDQL string = `
	SELECT *
	FROM license
	WHERE license_short_name = $1`

func (dao *LicenseDAO) GetLicenseByShortName(licenseShortName string) (*model.License, error) {
	var license model.License
	if err := dao.db.QueryRowx(getLicenseByShortNameDQL, strings.ToUpper(licenseShortName)).StructScan(&license); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("license with short name %s not found", strings.ToUpper(licenseShortName))
		}
		return nil, err
	}

	return &license, nil
}

// GetAllLicenses retrieves all categories from the database.
const getAllLicensesDQL string = `
	SELECT *
	FROM license`

func (dao *LicenseDAO) GetAllLicenses(ctx context.Context) ([]model.License, error) {
	var licenses []model.License
	if err := dao.db.SelectContext(ctx, &licenses, getAllLicensesDQL); err != nil {
		return nil, err
	}

	return licenses, nil
}

// CreateLicense inserts a new license into the database.
// Returns an error if the insertion fails.
// Does not return the PK as type tables have PK specified by the request.
func (dao *LicenseDAO) CreateLicense(licenseReq *model.LicenseRequest) error {
	_, err := dao.db.Exec(createLicenseDML,
		strings.ToUpper(licenseReq.LicenseShortName), licenseReq.LicenseFullName, licenseReq.LicenseUrl)
	if err != nil {
		return err
	}

	return nil
}

// UpdateLicense updates an existing license in the database.
const updateLicenseDML string = `
	UPDATE license
	SET
		license_full_name = $1,
		url = $2
	WHERE license_short_name = $3`

func (dao *LicenseDAO) UpdateLicense(licenseReq *model.LicenseRequest) error {
	result, err := dao.db.Exec(updateLicenseDML,
		licenseReq.LicenseFullName, licenseReq.LicenseUrl, strings.ToUpper(licenseReq.LicenseShortName))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("license with Short Name '%s' not found", licenseReq.LicenseShortName)
	}

	return nil
}

// DeleteLicense deletes a license from the database.
const deleteLicenseDML string = `
	DELETE FROM license
	WHERE license_short_name = $1`

func (dao *LicenseDAO) DeleteLicense(shortName string) error {
	result, err := dao.db.Exec(deleteLicenseDML, strings.ToUpper(shortName))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("license with Short Name '%s' not found", strings.ToUpper(shortName))
	}

	return nil
}
