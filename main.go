// Fluxy server main package
package main

import (
	"log"
	"net"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/thegeekywanderer/fluxy/config"
	database "github.com/thegeekywanderer/fluxy/db"
	"github.com/thegeekywanderer/fluxy/logger"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
	handler "github.com/thegeekywanderer/fluxy/pkg/v1/handlers/grpc"
	"github.com/thegeekywanderer/fluxy/pkg/v1/repository"
	"github.com/thegeekywanderer/fluxy/pkg/v1/usecase"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)


func main() {
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	viper.SetDefault("LOG_LEVEL", "DEBUG")
	logLevel := viper.GetString("LOG_LEVEL")

	if err := config.SetupConfig(); err != nil {
		logger.Error("%v", err)
	}
	logger.SetLogLevel(logLevel)

	db, err := database.DBConnection(config.GetDSNConfig())
	cache := database.RedisConnection(config.GetRedisConfig())
	if err != nil {
		logger.Fatal("%v", err)
	}
	database.Migrate(db)

	// add a listener address
	lis, err := net.Listen("tcp", config.ServerConfig())
	if err != nil {
		log.Fatalf("ERROR STARTING THE SERVER: %v", err)
	}

	// start the grpc server
	grpcServer := grpc.NewServer()

	clientUseCase := initClientServer(db.Database, cache)
	handler.NewServer(grpcServer, clientUseCase)

	// start serving to the address
	log.Fatal(grpcServer.Serve(lis))
}

func initClientServer(db *gorm.DB, cache *redis.Client) interfaces.UseCaseInterface{
	clientRepo := repository.New(db, cache)
	return usecase.New(clientRepo, viper.GetString("ALGORITHM"), cache)
}
