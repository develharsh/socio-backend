package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/harshvsinghme/socio-backend.git/dbUtils"
	globals "github.com/harshvsinghme/socio-backend.git/global"
	"github.com/harshvsinghme/socio-backend.git/router"

	"github.com/sirupsen/logrus"
)

func main() {

	InitLogrus()

	dbUtils.InitRedisConn()
	dbUtils.InitMongoConn()

	redisClient := dbUtils.RedisClient

	appRouter := router.InitRouter()
	//disconnect mongo, redis below when server fails
	defer func() {
		fmt.Println("Cleaning/Disconnecting Redis and MongoDB")
		// Cleanup Redis client
		if err := redisClient.Close(); err != nil {
			logrus.Error("Error closing Redis client:", err)
		}

		// Cleanup MongoDB client
		if err := dbUtils.MongoClient.Disconnect(context.Background()); err != nil {
			logrus.Error("Error disconnecting MongoDB client:", err)
		}
	}()

	logrus.Info("Running on port :3001")

	srv := &http.Server{
		Addr: ":3001",
		// Handler: appRouter,
		Handler: handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(appRouter),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Error(err.Error())
			logrus.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logrus.Error("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Error("Server Shutdown:")
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")

}

func InitLogrus() {

	fmt.Println("Initiating logger")

	Formatter := new(logrus.TextFormatter)
	Formatter.FullTimestamp = true
	Formatter.TimestampFormat = time.RFC3339
	logrus.SetFormatter(Formatter)

	// Add the UTC timestamp hook to logrus
	logrus.AddHook(UTCTimeHook{})

	globals.LoadGlobals(".env")

	fmt.Println("Running in " + globals.ENV + " mode")
}

// UTCTimeHook is a custom logrus hook that ensures UTC timestamps.
type UTCTimeHook struct{}

// Levels returns the log levels for which this hook should be called.
func (hook UTCTimeHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log event is fired.
func (hook UTCTimeHook) Fire(entry *logrus.Entry) error {
	entry.Time = entry.Time.UTC()
	return nil
}
