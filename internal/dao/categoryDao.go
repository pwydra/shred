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

type CategoryDaoInterface interface {
	CreateCategory(catReq *model.CategoryRequest) (model.Category, error)
	GetCategoryByCode(code string) (*model.Category, error)
	UpdateCategory(catReq *model.CategoryRequest) error
	DeleteCategory(code string) error
}

// Ensure ExerciseDao implements ExerciseDaoInterface
var _ CategoryDaoInterface = (*CategoryDAO)(nil)

// CategoryDAO provides access to the categories in the database.
type CategoryDAO struct {
	db *sqlx.DB
}

// NewCategoryDAO creates a new instance of CategoryDAO.
func NewCategoryDAO(db *sqlx.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

const createCatDML string = `
	INSERT INTO category_type (
		category_code, category_name, category_description, created_by
	) VALUES (
		$1, $2, $3, $4
	) RETURNING created_at, updated_at`

// GetCategoryByCode retrieves a category by its Code.
const getCatByCodeDQL string = `
	SELECT *
	FROM category_type
	WHERE category_code = $1`

func (dao *CategoryDAO) GetCategoryByCode(catCode string) (*model.Category, error) {
	var category model.Category
	if err := dao.db.QueryRowx(getCatByCodeDQL, strings.ToUpper(catCode)).StructScan(&category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category with code %s not found", strings.ToUpper(catCode))
		}
		return nil, err
	}

	return &category, nil
}

// GetAllCategories retrieves all categories from the database.
const getAllCatsDQL string = `
	SELECT *
	FROM category_type`

func (dao *CategoryDAO) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	if err := dao.db.SelectContext(ctx, &categories, getAllCatsDQL); err != nil {
		return nil, err
	}

	return categories, nil
}

// CreateCategory inserts a new category into the database.
// Returns an error if the insertion fails.
// Does not return the PK as type tables have PK specified by the request.
func (dao *CategoryDAO) CreateCategory(catReq *model.CategoryRequest) (model.Category, error) {
	cat := model.Category{
		CategoryFields: catReq.CategoryFields,
		AuditRecord:    model.AuditRecord{CreatedBy: catReq.CreatedBy},
	}
	err := dao.db.QueryRowx(createCatDML,
		strings.ToUpper(catReq.CategoryCode), catReq.CategoryName,
		catReq.CategoryDesc, catReq.CreatedBy).Scan(&cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

// UpdateCategory updates an existing category in the database.
const updateCatDML string = `
	UPDATE category_type
	SET
		category_name = $1,
		category_description = $2
	WHERE category_code = $3`

func (dao *CategoryDAO) UpdateCategory(catReq *model.CategoryRequest) error {
	result, err := dao.db.Exec(updateCatDML,
		catReq.CategoryName, catReq.CategoryDesc, strings.ToUpper(catReq.CategoryCode))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with Code %s not found", catReq.CategoryCode)
	}

	return nil
}

// DeleteCategory deletes a category from the database.
const deleteCatDML string = `
	DELETE FROM category_type
	WHERE category_code = $1`

func (dao *CategoryDAO) DeleteCategory(code string) error {
	result, err := dao.db.Exec(deleteCatDML, strings.ToUpper(code))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with code %s not found", strings.ToUpper(code))
	}

	return nil
}
