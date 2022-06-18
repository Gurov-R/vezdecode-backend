package main

import (
	vezdecodebackend "Gurov-R/vezdecode-backend"
	"Gurov-R/vezdecode-backend/pkg/cli"
	"Gurov-R/vezdecode-backend/pkg/handler"
	"Gurov-R/vezdecode-backend/pkg/repository"
	"Gurov-R/vezdecode-backend/pkg/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		gin.SetMode(gin.ReleaseMode)
		go func() {
			cli.RunCli()
		}()
	}

	server := new(vezdecodebackend.Server)
	if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
