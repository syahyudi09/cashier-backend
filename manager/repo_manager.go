package manager

import (
	"cashier/repository"
	"sync"
)

type RepositoryManager interface {
	GetUserRepo() repository.UserRepository
	GetCategoryRepo() repository.CategoryRepository
	GetProductRepo() repository.ProductRepository
}

type repositoryManager struct {
	infra        InfraManager
	userRepo     repository.UserRepository
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

var onceLoadUserRepo sync.Once
var onceLoadCategoryRepo sync.Once
var onceLoadProductRepo sync.Once

func (rm *repositoryManager) GetUserRepo() repository.UserRepository {
	onceLoadUserRepo.Do(func() {
		rm.userRepo = repository.NewUserRepository(rm.infra.GetDB())
	})
	return rm.userRepo
}

func (rm *repositoryManager) GetCategoryRepo() repository.CategoryRepository {
	onceLoadCategoryRepo.Do(func() {
		rm.categoryRepo = repository.NewCategoryRepository(rm.infra.GetDB())
	})
	return rm.categoryRepo
}

func (rm *repositoryManager) GetProductRepo() repository.ProductRepository {
	onceLoadProductRepo.Do(func() {
		rm.productRepo = repository.NewProductRepository(rm.infra.GetDB())
	})
	return rm.productRepo
}

func NewRepoManager(i InfraManager) RepositoryManager {
	return &repositoryManager{
		infra: i,
	}
}
