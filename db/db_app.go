package db

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres_user dbname=postgres_db password=postgres_password sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Fatal().Err(err).Msg("Database migration failed")
	}
}

type User struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Nickname  string `gorm:"size:255"`
	Email     string `gorm:"type:varchar(100);unique"`
	Password  string `gorm:"type:varchar(100)"`
}
