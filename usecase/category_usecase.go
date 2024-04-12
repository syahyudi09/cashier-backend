package usecase

import (
	"cashier/model"
	"cashier/repository"
	"cashier/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CategoryUsecase interface {
	CreateCategory(input model.CategoryInput) error
	CategoryNameExits(name string) (bool, error)
	FindById(id string) (*model.CategoryModel, error)
	Update(id string, updatedCategory model.CategoryUpdate) error
	GetAll(c context.Context, page, pageSize int, search string) ([]*model.CategoryModel, int, error)
	CheckNameForUpdate(name, id string) (bool, error)
	Delete(userID string) error
}

type categoryUsecase struct {
	categoryRepository repository.CategoryRepository
}

func (cr *categoryUsecase) CreateCategory(input model.CategoryInput) error {
	category := model.CategoryModel{}
	category.Id = utils.UuidGenerate()
	category.CategoryName = input.CategoryName
	category.Status = "Active"
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	return cr.categoryRepository.Create(category)
}

func (cr *categoryUsecase) CategoryNameExits(name string) (bool, error) {
	return cr.categoryRepository.CheckCategoryName(name)
}

func (cr *categoryUsecase) Update(id string, updatedCategory model.CategoryUpdate) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidUUIDFormat
	}

	category, err := cr.categoryRepository.FindById(id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}

	if updatedCategory.CategoryName != "" {
		category.CategoryName = updatedCategory.CategoryName
	}

	if updatedCategory.Status != "" {
		category.Status = updatedCategory.Status
	}

	err = cr.categoryRepository.Update(id, *category)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (cr *categoryUsecase) GetAll(c context.Context, page, pageSize int, search string) ([]*model.CategoryModel, int, error) {
	category, totalDocs, err := cr.categoryRepository.GetAll(c, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching category: %v", err)
	}

	return category, totalDocs, nil
}

func (cr *categoryUsecase) FindById(id string) (*model.CategoryModel, error) {
	return cr.categoryRepository.FindById(id)
}

func (cr *categoryUsecase) CheckNameForUpdate(name, id string) (bool, error) {
	exists, err := cr.categoryRepository.CheckNameForUpdate(name, id)
	if err != nil {
		return false, fmt.Errorf("error checking name for update: %w", err)
	}
	return exists, nil
}

func (cr *categoryUsecase) Delete(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidUUIDFormat
	}
	return cr.categoryRepository.Delete(userID)
}

func NewCategoryUsecase(cr repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepository: cr,
	}
}
