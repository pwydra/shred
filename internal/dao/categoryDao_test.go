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

func TestGetCategoryByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catCode := "STRENGTH"
	mock.ExpectQuery("SELECT \\* FROM category_type WHERE category_code = \\$1").
		WithArgs(catCode).
		WillReturnRows(sqlmock.NewRows([]string{"category_code", "category_name", "category_description"}).
			AddRow(catCode, "Strength", "Strength training exercises"))

	category, err := dao.GetCategoryByCode(catCode)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Strength", category.CategoryName)
}

func TestGetCategoryByCode_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catCode := "INVALID"
	mock.ExpectQuery("SELECT \\* FROM category_type WHERE category_code = \\$1").
		WithArgs(catCode).
		WillReturnError(sql.ErrNoRows)

	category, err := dao.GetCategoryByCode(catCode)
	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Equal(t, "category with code INVALID not found", err.Error())
}

func TestGetAllCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	mock.ExpectQuery("SELECT \\* FROM category_type").
		WillReturnRows(sqlmock.NewRows([]string{"category_code", "category_name", "category_description"}).
			AddRow("STRENGTH", "Strength", "Strength training exercises").
			AddRow("CARDIO", "Cardio", "Cardiovascular exercises"))

	ctx := context.Background()
	categories, err := dao.GetAllCategories(ctx)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Strength", categories[0].CategoryName)
	assert.Equal(t, "Cardio", categories[1].CategoryName)
}

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.CategoryRequest{
		CategoryFields: model.CategoryFields{
			CategoryCode: "STRENGTH",
			CategoryName: "Strength",
			CategoryDesc: "Strength training exercises",
		},
	}

	mock.ExpectExec("INSERT INTO category_type \\( category_code, category_name, category_description \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(catReq.CategoryCode, catReq.CategoryName, catReq.CategoryDesc).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.CreateCategory(catReq)
	assert.NoError(t, err)
}

func TestCreateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.CategoryRequest{
		CategoryFields: model.CategoryFields{
			CategoryCode: "STRENGTH",
			CategoryName: "Strength",
			CategoryDesc: "Strength training exercises",
		},
	}

	mock.ExpectExec("INSERT INTO category_type \\( category_code, category_name, category_description \\) VALUES \\( \\$1, \\$2, \\$3 \\)").
		WithArgs(catReq.CategoryCode, catReq.CategoryName, catReq.CategoryDesc).
		WillReturnError(errors.New("insertion error"))

	err = dao.CreateCategory(catReq)
	assert.Error(t, err)
	assert.Equal(t, "insertion error", err.Error())
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.CategoryRequest{
		CategoryFields: model.CategoryFields{
			CategoryCode: "STRENGTH",
			CategoryName: "Strength",
			CategoryDesc: "Strength training exercises",
		},
	}

	mock.ExpectExec("UPDATE category_type SET category_name = \\$1, category_description = \\$2 WHERE category_code = \\$3").
		WithArgs(catReq.CategoryName, catReq.CategoryDesc, catReq.CategoryCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.UpdateCategory(catReq)
	assert.NoError(t, err)
}

func TestUpdateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.CategoryRequest{
		CategoryFields: model.CategoryFields{
			CategoryCode: "STRENGTH",
			CategoryName: "Strength",
			CategoryDesc: "Strength training exercises",
		},
	}

	mock.ExpectExec("UPDATE").
		WithArgs(catReq.CategoryName, catReq.CategoryDesc, catReq.CategoryCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.UpdateCategory(catReq)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestUpdateCategory_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catReq := &model.CategoryRequest{
		CategoryFields: model.CategoryFields{
			CategoryCode: "INVALID",
			CategoryName: "Invalid",
			CategoryDesc: "Invalid description",
		},
	}

	mock.ExpectExec("UPDATE category_type SET category_name = \\$1, category_description = \\$2 WHERE category_code = \\$3").
		WithArgs(catReq.CategoryName, catReq.CategoryDesc, catReq.CategoryCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.UpdateCategory(catReq)
	assert.Error(t, err)
	assert.Equal(t, "category with Code INVALID not found", err.Error())
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catCode := "STRENGTH"
	mock.ExpectExec("DELETE FROM category_type WHERE category_code = \\$1").
		WithArgs(catCode).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.DeleteCategory(catCode)
	assert.NoError(t, err)
}

func TestDeleteCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catCode := "STRENGTH"
	mock.ExpectExec("DELETE FROM category_type WHERE category_code = \\$1").
		WithArgs(catCode).
		WillReturnError(sqlmock.ErrCancelled)

	err = dao.DeleteCategory(catCode)
	assert.Error(t, err)
	assert.Equal(t, "canceling query due to user request", err.Error())
}

func TestDeleteCategory_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dao := NewCategoryDAO(sqlx.NewDb(db, "postgres"))

	catCode := "INVALID"
	mock.ExpectExec("DELETE FROM category_type WHERE category_code = \\$1").
		WithArgs(catCode).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = dao.DeleteCategory(catCode)
	assert.Error(t, err)
	assert.Equal(t, "category with code INVALID not found", err.Error())
}
