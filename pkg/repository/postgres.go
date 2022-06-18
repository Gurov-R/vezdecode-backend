package repository

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	memesTable = "memes"
	likesTable = "likes"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	if len(os.Args) > 1 && os.Args[1] == "db-up" {
		schema, err := ioutil.ReadFile("./schema/init.sql")
		if err != nil {
			logrus.Fatalf("Error while reading ./schema/init.sql")
		}

		db.MustExec(string(schema))

		fmt.Println("Migration complete.")
		os.Exit(0)
	}
	if len(os.Args) > 1 && os.Args[1] == "db-down" {
		schema, err := ioutil.ReadFile("./schema/down.sql")
		if err != nil {
			logrus.Fatalf("Error while reading ./schema/down.sql")
		}

		db.MustExec(string(schema))

		fmt.Println("Migration complete.")
		os.Exit(0)
	}

	return db, nil
}

