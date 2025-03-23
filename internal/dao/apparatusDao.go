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

// ApparatusDAO provides access to the apparatuses in the database.
type ApparatusDAO struct {
	db *sqlx.DB
}

// NewApparatusDAO creates a new instance of ApparatusDAO.
func NewApparatusDAO(db *sqlx.DB) *ApparatusDAO {
	return &ApparatusDAO{db: db}
}

const createAppDML string = `
	INSERT INTO apparatus_type (
		apparatus_code, apparatus_name, apparatus_description
	) VALUES (
		$1, $2, $3
	)`

// GetApparatusByCode retrieves a apparatus by its Code.
const getAppByCodeDQL string = `
	SELECT *
	FROM apparatus_type
	WHERE apparatus_code = $1`
func (dao *ApparatusDAO) GetApparatusByCode(appCode string) (*model.Apparatus, error) {
	var apparatus model.Apparatus
	if err := dao.db.QueryRowx(getAppByCodeDQL, strings.ToUpper(appCode)).StructScan(&apparatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("apparatus with code %s not found", strings.ToUpper(appCode))
		}
		return nil, err
	}

	return &apparatus, nil
}

// GetAllApparatuses retrieves all categories from the database.
const getAllAppsDQL string = `
	SELECT *
	FROM apparatus_type`
func (dao *ApparatusDAO) GetAllApparatuses(ctx context.Context) ([]model.Apparatus, error) {
	var apparatuses []model.Apparatus
	if err := dao.db.SelectContext(ctx, &apparatuses, getAllAppsDQL); err != nil {
		return nil, err
	}

	return apparatuses, nil
}

// CreateApparatus inserts a new apparatus into the database.
// Returns an error if the insertion fails.
// Does not return the PK as type tables have PK specified by the request.
func (dao *ApparatusDAO) CreateApparatus(appReq *model.ApparatusRequest) (error) {
	_, err := dao.db.Exec(createAppDML,
		strings.ToUpper(appReq.ApparatusCode), appReq.ApparatusName, appReq.ApparatusDesc)
	if err != nil {
		return err
	}

	return nil
}

// UpdateApparatus updates an existing apparatus in the database.
const updateAppDML string = `
	UPDATE apparatus_type
	SET
		apparatus_name = $1,
		apparatus_description = $2
	WHERE apparatus_code = $3`
func (dao *ApparatusDAO) UpdateApparatus(appReq *model.ApparatusRequest) error {
	result, err := dao.db.Exec(updateAppDML, 
		appReq.ApparatusName, appReq.ApparatusDesc, strings.ToUpper(appReq.ApparatusCode))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("apparatus with Code %s not found", appReq.ApparatusCode)
	}

	return nil
}

// DeleteApparatus deletes a apparatus from the database.
const deleteAppDML string = `
	DELETE FROM apparatus_type
	WHERE apparatus_code = $1`

func (dao *ApparatusDAO) DeleteApparatus(code string) error {
	result, err := dao.db.Exec(deleteAppDML, strings.ToUpper(code))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("apparatus with code %s not found", strings.ToUpper(code))
	}

	return nil
}
