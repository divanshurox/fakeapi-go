package server

import (
	"FakeAPI/internal/db"
	"FakeAPI/internal/logger"
	"FakeAPI/internal/middleware"
	"FakeAPI/internal/mongo"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type GracefulServer struct {
	listener      net.Listener
	httpServer    *http.Server
	hystrixServer *hystrix.StreamHandler
}

func NewServer(port string) *GracefulServer {
	mux := chi.NewRouter()
	mux.Use(middleware.AddLogging)
	mux.Use(middleware.AddRateLimiting)
	mux.Route("/v1", RouteFakeApi)
	httpServer := &http.Server{Addr: ":" + port, Handler: mux}
	return &GracefulServer{httpServer: httpServer}
}

func (gc *GracefulServer) Prestart() error {
	logger.InitLogger()
	if logger.GetLogger() == nil {
		errMsg := "failed to initialize logger"
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	_, err := db.GetDatabase(mongo.GetInstance()).Connect(
		db.NewEmptyConfig().WithHost("localhost").WithPort("27017").WithUsername("root").WithPassword("example"),
	)
	if err != nil {
		logger.GetLogger().Error("Unable to connect to MongoDB", zap.Error(err))
		return err
	}

	gc.hystrixServer = hystrix.NewStreamHandler()
	if gc.hystrixServer != nil {
		gc.hystrixServer.Start()
		go func() {
			err := http.ListenAndServe(net.JoinHostPort("", "81"), gc.hystrixServer)
			if err != nil {
				logger.GetLogger().Error("error with hystrix", zap.Error(err))
				gc.hystrixServer.Stop()
			}
		}()
	}
	return nil
}

func (gc *GracefulServer) Start() (chan bool, error) {
	listener, err := net.Listen("tcp", gc.httpServer.Addr)
	if err != nil {
		return nil, err
	}
	gc.listener = listener
	go gc.httpServer.Serve(gc.listener)
	done := make(chan bool)
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-interrupt
		logger.GetLogger().Info("Signal intercepted", zap.String("Signal", sig.String()))
		gc.Stop()
		done <- true
	}()
	return done, nil
}

func (gc *GracefulServer) Stop() error {
	logger.Close()
	if gc.listener != nil {
		err := gc.listener.Close()
		if err != nil {
			return err
		}
		gc.listener = nil
	}
	if mongo.GetInstance() != nil {
		db.GetDatabase(mongo.GetInstance()).Close()
	}
	if gc.hystrixServer != nil {
		gc.hystrixServer.Stop()
		gc.hystrixServer = nil
	}
	logger.GetLogger().Info("Shutting down server")
	return nil
}
