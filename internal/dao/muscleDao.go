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

// MuscleDAO provides access to the muscles in the database.
type MuscleDAO struct {
	db *sqlx.DB
}

// NewMuscleDAO creates a new instance of MuscleDAO.
func NewMuscleDAO(db *sqlx.DB) *MuscleDAO {
	return &MuscleDAO{db: db}
}

const createMusDML string = `
	INSERT INTO muscle_type (
		muscle_code, muscle_name, muscle_description, muscle_group
	) VALUES (
		$1, $2, $3, $4
	)`

// GetMuscleByCode retrieves a muscle by its Code.
const getMusByCodeDQL string = `
	SELECT *
	FROM muscle_type
	WHERE muscle_code = $1`
func (dao *MuscleDAO) GetMuscleByCode(musCode string) (*model.Muscle, error) {
	var muscle model.Muscle
	if err := dao.db.QueryRowx(getMusByCodeDQL, strings.ToUpper(musCode)).StructScan(&muscle); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("muscle with code %s not found", strings.ToUpper(musCode))
		}
		return nil, err
	}

	return &muscle, nil
}

// GetAllMuscles retrieves all muscles from the database.
const getAllMusDQL string = `
	SELECT *
	FROM muscle_type`
func (dao *MuscleDAO) GetAllMuscles(ctx context.Context) ([]model.Muscle, error) {
	var muscles []model.Muscle
	if err := dao.db.SelectContext(ctx, &muscles, getAllMusDQL); err != nil {
		return nil, err
	}

	return muscles, nil
}

// CreateMuscle inserts a new muscle into the database.
// Returns an error if the insertion fails.
// Does not return the PK as type tables have PK specified by the request.
func (dao *MuscleDAO) CreateMuscle(musReq *model.MuscleRequest) (error) {
	_, err := dao.db.Exec(createMusDML,
		strings.ToUpper(musReq.MuscleCode), musReq.MuscleName, musReq.MuscleDesc)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMuscle updates an existing muscle in the database.
const updateMusDML string = `
	UPDATE muscle_type
	SET
		muscle_name = $1,
		muscle_description = $2,
		muscle_group = $3
	WHERE muscle_code = $4`
func (dao *MuscleDAO) UpdateMuscle(musReq *model.MuscleRequest) error {
	result, err := dao.db.Exec(updateMusDML, 
		musReq.MuscleName, musReq.MuscleDesc, musReq.MuscleGroup, strings.ToUpper(musReq.MuscleCode))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("muscle with Code %s not found", musReq.MuscleCode)
	}

	return nil
}

// DeleteMuscle deletes a muscle from the database.
const deleteMusDML string = `
	DELETE FROM muscle_type
	WHERE muscle_code = $1`

func (dao *MuscleDAO) DeleteMuscle(code string) error {
	result, err := dao.db.Exec(deleteMusDML, strings.ToUpper(code))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("muscle with code %s not found", strings.ToUpper(code))
	}

	return nil
}
