package main

import (
	"avitoTech"
	"avitoTech/pkg/handler"
	"avitoTech/pkg/repository"
	"avitoTech/pkg/service"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// инициализация бд
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db %s", err.Error())
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	srv := new(avitoTech.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("error: main: occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("avitoTech started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("avitoTech shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Error("error: main: occured on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Error("error: main: occured on server db connection closed: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
