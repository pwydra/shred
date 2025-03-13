package dao

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pwydra/shred/internal/model"
)

type ExerciseDao struct {
	db *sql.DB
}

func NewExerciseDao(db *sql.DB) *ExerciseDao {
	return &ExerciseDao{db: db}
}

/*
	insert into exercise (
	   exercise_name, description, instructions,
	   category, cues, primary_muscles,
	   secondary_muscles, front_image, back_image,
	   video_url, apparatus, license, user_uuid

) values ('Squat', 'Lower Body', 'Stand with feet shoulder-width apart, toes slightly turned out. Push hips back and bend knees to lower into a squat. Keep chest up and back flat. Push through heels to stand back up.',

	'Strength', 'Keep chest up and back flat; imagine pushing through floor; knees travel over toes; slow eccentric', 'Quadriceps',
	'Glutes', null, null,
	null, 'Barbell', 'None',
	'f47f3b9e-0b1d-4b7b-8b3d-3b1b1b1b1b1b');
*/
const createDML string = `
	INSERT INTO exercise (
		exercise_name, description, cues, primary_muscles, apparatus, user_uuid
	) VALUES ($1, $2, $3, $4, $5, $6) RETURNING uuid`
const updateDML string = "UPDATE exercise SET exercise_name = $1, description = $2, cues = $3, primary_muscles = $4, apparatus = $5 WHERE uuid = $6"
const deleteDML string = "DELETE FROM exercise WHERE uuid = $1"
const request1DQL string = `
	SELECT uuid, exercise_name, description, cues, primary_muscles, apparatus, created_at, user_uuid 
	FROM   exercise 
	WHERE  uuid = $1`

//var requestAllDQL string = "SELECT uuid, exercise_name, description cues, primary_muscles, apparatus, created_at, user_uuid FROM exercises"

func (dao *ExerciseDao) Create(exerciseRequest *model.ExerciseRequest) (*model.Exercise, error) {
	exercise := model.Exercise{
		ExerciseFields: exerciseRequest.ExerciseFields,
	}

	err := dao.db.QueryRow(createDML,
		exerciseRequest.Name, exerciseRequest.Description, exerciseRequest.Cues, exerciseRequest.PrimaryMuscles,
		exerciseRequest.Apparatus, exerciseRequest.UserUuid).Scan(&exercise.Uuid)
	if err != nil {
		log.Println("Error creating exercise:", err)
		return nil, err
	}
	return &exercise, nil
}

func (dao *ExerciseDao) Read(uuid uuid.UUID) (*model.Exercise, error) {
	fmt.Println("Reading exercise with uuid:", uuid)
	row := dao.db.QueryRow(request1DQL, uuid)
	exercise := &model.Exercise{}
	fmt.Println("scanning exercise from row", row)
	err := row.Scan(
		&exercise.Uuid, &exercise.Name, &exercise.Description, &exercise.Cues,
		&exercise.PrimaryMuscles, &exercise.Apparatus, &exercise.CreateDt, &exercise.UserUuid)
	if err != nil {
		log.Println("Error reading exercise:", err)
		return nil, err
	}
	return exercise, nil
}

func (dao *ExerciseDao) Update(exercise *model.Exercise) error {
	_, err := dao.db.Exec(updateDML, exercise.Name, exercise.Description, exercise.Cues, exercise.PrimaryMuscles, exercise.Apparatus, exercise.Uuid)
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
