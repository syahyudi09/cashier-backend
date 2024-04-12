package manager

import (
	"cashier/usecase"
	"sync"
)

type UseceaseManager interface {
	GetUserUsecase() usecase.UserUsecase
	GetCategoryUsecase() usecase.CategoryUsecase
	GetProductUsecase() usecase.ProductUsecase
}

type usecaseManager struct {
	rm              RepositoryManager
	userUsecase     usecase.UserUsecase
	categoryUsecase usecase.CategoryUsecase
	productUsecase  usecase.ProductUsecase
}

var onceLoadUserUsecase sync.Once
var onceLoadCategoryUsecase sync.Once
var onceLoadProductUsecase sync.Once

func (um *usecaseManager) GetUserUsecase() usecase.UserUsecase {
	onceLoadUserUsecase.Do(func() {
		um.userUsecase = usecase.NewUserUsecase(
			um.rm.GetUserRepo(),
		)
	})
	return um.userUsecase
}

func (um *usecaseManager) GetCategoryUsecase() usecase.CategoryUsecase {
	onceLoadCategoryUsecase.Do(func() {
		um.categoryUsecase = usecase.NewCategoryUsecase(
			um.rm.GetCategoryRepo(),
		)
	})
	return um.categoryUsecase
}

func (um *usecaseManager) GetProductUsecase() usecase.ProductUsecase {
	onceLoadProductUsecase.Do(func() {
		um.productUsecase = usecase.NewProductUsecase(
			um.rm.GetProductRepo(),
		)
	})
	return um.productUsecase
}

func NewUsecasemanager(rm RepositoryManager) UseceaseManager {
	return &usecaseManager{
		rm: rm,
	}
}
