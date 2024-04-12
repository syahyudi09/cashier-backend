package delivery

import (
	"cashier/config"
	"cashier/delivery/controller"
	"cashier/manager"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Run()
}

type serverImpl struct {
	engine  *gin.Engine
	usecase manager.UseceaseManager
	config  config.Config
}

func (s *serverImpl) Run() {
	controller.NewUserController(s.engine, s.usecase.GetUserUsecase())
	controller.NewCategoryController(s.engine, s.usecase.GetCategoryUsecase())
	controller.NewProductController(s.engine, s.usecase.GetProductUsecase())
	controller.NewAuthController(s.engine, s.usecase.GetUserUsecase())

	s.engine.Run(":8080")
}

func NewServer() Server {
	config, err := config.NewConfig()
	if err != nil{
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // Sesuaikan dengan origin Anda
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"}, // Tambahkan 'Authorization' ke header yang diizinkan
	}))

	infra := manager.NewInfraManager(config)
	repo := manager.NewRepoManager(infra)
	usecase := manager.NewUsecasemanager(repo)

	return &serverImpl{
		engine:  r,
		usecase: usecase,
		config: config,
	}
}
