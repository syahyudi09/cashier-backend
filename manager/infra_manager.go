package manager

import (
	"database/sql"
	"fmt"
	"cashier/config"
	"sync"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type InfraManager interface {
	GetDB() *sql.DB
}

type infraManagerImpl struct {
	db  *sql.DB
	cfg config.Config
}

var onceLoadDB sync.Once

func (i *infraManagerImpl) InitDB() {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", i.cfg.Host, i.cfg.Port, i.cfg.User, i.cfg.Password, i.cfg.Name)
	var err error
	onceLoadDB.Do(func() {
		i.db, err = sql.Open("postgres", psqlConn)
		if err != nil {
			panic(err)
		}
	})
}

func (i *infraManagerImpl) GetDB() *sql.DB {
	return i.db
}

func NewInfraManager(config config.Config) InfraManager {
	infra := &infraManagerImpl{
		cfg: config,
	}
	infra.InitDB()
	return infra
}