package repository

import (
	"cashier/model"
	"cashier/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CategoryRepository interface {
	Create(newCategory model.CategoryModel) error
	CheckCategoryName(name string) (bool, error)
	FindById(id string) (*model.CategoryModel, error)
	Update(id string, updatedCategory model.CategoryModel) error
	GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.CategoryModel, int, error)
	CheckNameForUpdate(nama, id string) (bool, error)
	Delete(ID string) error
}

type categoryRepository struct {
	db *sql.DB
}

func (cr *categoryRepository) Create(newCategory model.CategoryModel) error {
	newCategory.Id = utils.UuidGenerate()
	insertQuery := "INSERT INTO categories(id, category_name, status, created_at, updated_at) VALUES($1,$2,$3,$4, $5) RETURNING id"

	err := cr.db.QueryRow(insertQuery, newCategory.Id, newCategory.CategoryName, newCategory.Status, newCategory.CreatedAt, newCategory.UpdatedAt).Scan(&newCategory.Id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to create category: %v", err)
	}
	return nil
}

func (cr *categoryRepository) Update(id string, updatedCategory model.CategoryModel) error {
	updateQuery := "UPDATE categories SET category_name = $1, status = $2, updated_at = $3 WHERE id = $4"

	_, err := cr.db.Exec(updateQuery, updatedCategory.CategoryName, updatedCategory.Status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}

	return nil
}

func (cr *categoryRepository) FindById(id string) (*model.CategoryModel, error) {
	getQuery := "SELECT id, category_name, status, created_at, updated_at FROM categories WHERE id = $1"

	row := cr.db.QueryRow(getQuery, id)

	category := &model.CategoryModel{}
	err := row.Scan(&category.Id, &category.CategoryName, &category.Status, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, fmt.Errorf("error fetching user by email: %w", err)
	}

	return category, nil
}

func (cr *categoryRepository) GetAll(ctx context.Context, page, pageSize int, search string) ([]*model.CategoryModel, int, error) {
	offset := (page - 1) * pageSize
	getQuery := `
		SELECT id, category_name, status, created_at, updated_at FROM categories WHERE category_name LIKE '%' || $1 || '%'
		ORDER BY id 
		LIMIT $2 
		OFFSET $3`

	rows, err := cr.db.QueryContext(ctx, getQuery, search, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error on categoryRepository.GetAll: %w", err)
	}
	defer rows.Close()

	var arrCategory []*model.CategoryModel

	for rows.Next() {
		category := &model.CategoryModel{}
		if err := rows.Scan(
			&category.Id, &category.CategoryName, &category.Status, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("error scanning campaign row: %w", err)
		}
		arrCategory = append(arrCategory, category)
	}

	countQuery := "SELECT COUNT(id) FROM categories" // Menyesuaikan query COUNT dengan offset
	var totalDocs int
	if err := cr.db.QueryRowContext(ctx, countQuery).Scan(&totalDocs); err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %w", err)
	}

	return arrCategory, totalDocs, nil
}

func (cr *categoryRepository) CheckCategoryName(name string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM categories WHERE category_name = $1)"

	var exists bool
	err := cr.db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error fetching user by category name: %w", err)
	}

	return exists, nil
}

func (cr *categoryRepository) CheckNameForUpdate(nama, id string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM categories WHERE category_name = $1 AND id != $2)"

	var exists bool
	err := cr.db.QueryRow(query, nama, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking existing ctaegory name for update: %w", err)
	}

	return exists, nil
}

func (cr *categoryRepository) Delete(ID string) error {
	deleteQuery := "DELETE FROM categories WHERE id = $1"

	_, err := cr.db.Query(deleteQuery, ID)
	if err != nil {
		return fmt.Errorf("error fetching categories by Id: %w", err)
	}
	return nil
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
