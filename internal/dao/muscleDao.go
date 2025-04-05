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

type MuscleDaoInterface interface {
	CreateMuscle(catReq *model.MuscleRequest) (model.Muscle, error)
	GetMuscleByCode(code string) (*model.Muscle, error)
	GetAllMuscles(ctx context.Context) ([]model.Muscle, error)
	UpdateMuscle(catReq *model.MuscleRequest) error
	DeleteMuscle(code string) error
}

// Ensure ExerciseDao implements ExerciseDaoInterface
var _ MuscleDaoInterface = (*MuscleDAO)(nil)

// NewMuscleDAO creates a new instance of MuscleDAO.
func NewMuscleDAO(db *sqlx.DB) *MuscleDAO {
	return &MuscleDAO{db: db}
}

const createMusDML string = `
	INSERT INTO muscle_type (
		muscle_code, muscle_name, muscle_description, muscle_group, created_by
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING created_at, updated_at
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
// Returns created object.
func (dao *MuscleDAO) CreateMuscle(musReq *model.MuscleRequest) (model.Muscle, error) {
	musReq.MuscleCode = strings.ToUpper(musReq.MuscleCode)
	mus := model.Muscle{
		MuscleFields: musReq.MuscleFields,
		AuditRecord:  model.AuditRecord{CreatedBy: musReq.CreatedBy},
	}
	err := dao.db.QueryRowx(createMusDML,
		musReq.MuscleCode, musReq.MuscleName, musReq.MuscleDesc,
		musReq.MuscleGroup, musReq.CreatedBy).Scan(&mus.CreatedAt, &mus.UpdatedAt)
	if err != nil {
		return mus, err
	}

	return mus, nil
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
