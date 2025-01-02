package database

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"fashora-backend/config"
	"fashora-backend/logger"
)

var once sync.Once

var dbAdapter DBAdapter

// DBAdapter interface represent adapter connect to DB
type DBAdapter interface {
	Open(cfg config.DbPostGreSQLConfig) error
	DB() *gorm.DB
	Connection() *gorm.DB
}

type adapter struct {
	connection *gorm.DB
	session    *gorm.DB
}

// Open opens a DB connection.
func (db *adapter) Open(cfg config.DbPostGreSQLConfig) error {
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,       // Slow SQL threshold
			LogLevel:      gormLogger.Silent, // Log level
			Colorful:      false,             // Disable color
		},
	)

	DB, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port),
		),
		&gorm.Config{
			Logger: newLogger,
		},
	)

	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Errorf("error %#v", err)
		return err
	}

	logger.Infof("max life time: %s", cfg.ConnMaxLifetimeSec)
	logger.Infof("max open connections: %s", cfg.MaxOpenCons)
	logger.Infof("max open idle connections: %s", cfg.MaxIdleCons)

	maxOpenCons, _ := strconv.Atoi(cfg.MaxOpenCons)
	maxIdleCons, _ := strconv.Atoi(cfg.MaxIdleCons)
	connMaxLifeTimeSec, _ := strconv.Atoi(cfg.ConnMaxLifetimeSec)
	sqlDB.SetMaxOpenConns(maxOpenCons)
	sqlDB.SetMaxIdleConns(maxIdleCons)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifeTimeSec) * time.Minute)

	db.connection = DB
	db.connection.Exec("Set time_zone = '+00:00'")
	db.session = db.connection.Session(&gorm.Session{})
	return nil
}

func (db *adapter) DB() *gorm.DB {
	return db.session
}

func (db *adapter) Connection() *gorm.DB {
	return db.connection
}

// NewDB returns a new instance of DB.
func newDB() DBAdapter {
	return &adapter{}
}

func GetDBInstance() DBAdapter {
	if dbAdapter == nil {
		once.Do(func() {
			dbAdapter = newDB()
		})
	}
	return dbAdapter
}
