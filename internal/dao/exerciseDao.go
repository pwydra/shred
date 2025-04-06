package dao

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	"github.com/pwydra/shred/internal/model"
)

type ExerciseDao struct {
	db *sqlx.DB
}

type ExerciseDaoInterface interface {
	Create(exerciseRequest *model.ExerciseRequest) (*model.Exercise, error)
	Read(uuid uuid.UUID) (*model.Exercise, error)
	Update(exercise *model.Exercise) error
	Delete(uuid uuid.UUID) error
}

// Ensure ExerciseDao implements ExerciseDaoInterface
var _ ExerciseDaoInterface = (*ExerciseDao)(nil)

func NewExerciseDao(db *sqlx.DB) *ExerciseDao {
	return &ExerciseDao{db: db}
}

const createDML string = `
	INSERT INTO exercise (
		exercise_name, exercise_description, instructions, cues, 
		video_url, category_code, license_short_name, license_author, 
		created_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING uuid, created_at, updated_at`

const updateDML string = `
	UPDATE exercise SET 
		exercise_name = $1, 
		exercise_description = $2,
		instructions = $3, 
		cues = $4,
		video_url = $5,
		category_code = $6,
		license_short_name = $7,
		license_author = $8
	WHERE exercise_uuid = $9`

const deleteDML string = "DELETE FROM exercise WHERE exercise_uuid = $1"

const request1DQL string = `
	SELECT * 
	FROM   exercise 
	WHERE  exercise_uuid = $1`

// TODO: add query for all exercises
// var requestAllDQL string = "SELECT uuid, exercise_name, description cues, primary_muscles, apparatus, created_at, user_uuid FROM exercises"

func (dao *ExerciseDao) Create(exReq *model.ExerciseRequest) (*model.Exercise, error) {
	exercise := model.Exercise{
		ExerciseFields: exReq.ExerciseFields,
	}

	err := dao.db.QueryRowx(createDML,
		exReq.ExerciseName, exReq.Description, exReq.Instructions, exReq.Cues,
		exReq.VideoUrl, exReq.CategoryCode, exReq.LicenseShortName, exReq.LicenseAuthor,
		exReq.CreatedBy).Scan(&exercise.ExerciseUuid, &exercise.CreatedAt, &exercise.UpdatedAt)
	if err != nil {
		log.Println("Error creating exercise:", err)
		return nil, err
	}
	return &exercise, nil
}

func (dao *ExerciseDao) Read(exUuid uuid.UUID) (*model.Exercise, error) {
	var ex model.Exercise
	err := dao.db.QueryRowx(request1DQL, exUuid).StructScan(&ex)
	if err != nil {
		log.Println("Error reading exercise:", err)
		return nil, err
	}
	return &ex, nil
}

func (dao *ExerciseDao) Update(exercise *model.Exercise) error {
	_, err := dao.db.Exec(updateDML,
		exercise.ExerciseName, exercise.Description, exercise.Instructions, exercise.Cues,
		exercise.VideoUrl, exercise.CategoryCode, exercise.LicenseShortName,
		exercise.LicenseAuthor, exercise.ExerciseUuid)
	if err != nil {
		log.Println("Error updating exercise:", err)
		return err
	}
	// TODO: return updated exercise
	return nil
}

func (dao *ExerciseDao) Delete(uuid uuid.UUID) error {
	_, err := dao.db.Exec(deleteDML, uuid)
	if err != nil {
		log.Println("Error deleting exercise:", err)
		return err
	}
	return nil
}
