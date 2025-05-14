package grpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	pb "github.com/dinklen/GolangCalc_V2/api/proto/generated"
	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/evaluator"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/neterr"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.UnimplementedCalculatorServer
	workerPool *WorkerPool
	config     *config.Config
	logger     *zap.Logger
}

func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{
		config:     cfg,
		workerPool: NewWorkerPool(cfg.Microservice.ComputingPower, cfg, logger),
		logger:     logger,
	}
}

type WorkerPool struct {
	tasks    chan *pb.Subexpression
	results  chan *pb.Result
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewWorkerPool(quantity int, cfg *config.Config, logger *zap.Logger) *WorkerPool {
	pool := &WorkerPool{
		tasks:    make(chan *pb.Subexpression),
		results:  make(chan *pb.Result),
		shutdown: make(chan struct{}),
	}

	pool.wg.Add(quantity)
	for i := 0; i < quantity; i++ {
		go pool.worker(cfg, logger)
	}

	return pool
}

func (p *WorkerPool) worker(cfg *config.Config, logger *zap.Logger) {
	logger.Info("worker started")
	defer p.wg.Done()

	for {
		select {
		case task := <-p.tasks:
			logger.Info("worker got a new task")
			res, err := evaluateExpression(task.LeftValue, task.RightValue, task.Operator, cfg, logger)
			response := &pb.Result{
				Id:       task.Id,
				ParentId: task.ParentId,
			}

			if err != nil {
				response.Error = err.Error()
			} else {
				response.Value = res
			}
			logger.Info("worker completed a task")
			p.results <- response

		case <-p.shutdown:
			logger.Info("shutdown worker")
			return
		}
	}
}

func (p *WorkerPool) Shutdown(logger *zap.Logger) {
	logger.Info("shutting down workers")
	close(p.shutdown)
	p.wg.Wait()
	close(p.results)
}

func evaluateExpression(leftValue, rightValue float64, operator string, cfg *config.Config, logger *zap.Logger) (float64, error) {
	result, err := evaluator.Calculate(leftValue, rightValue, operator, cfg, logger)
	if err != nil {
		return 0.0, fmt.Errorf("%s", err.Error())
	}

	logger.Info("task evaluating success")
	return result, nil
}

func (s *Server) Calculate(stream pb.Calculator_CalculateServer) error {
	defer s.workerPool.Shutdown(s.logger)

	go func() {
		for result := range s.workerPool.results {
			s.logger.Info("server got a new task")
			if err := stream.Send(result); err != nil {
				s.logger.Error("failed to response the task")
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			s.logger.Error("client connection closed")
			return neterr.ErrGRPCServerConnectionClosed
		}

		select {
		case s.workerPool.tasks <- req:
			s.logger.Info("worker got a new request")
		default:
			s.workerPool.results <- &pb.Result{
				Id:    req.Id,
				Error: "server is busy",
			}
		}
	}
}

func Start(cfg *config.Config, logger *zap.Logger) error {
	creds, err := credentials.NewServerTLSFromFile("../../certs/cert.pem", "../../certs/key.pem")
	if err != nil {
		logger.Fatal(neterr.ErrGRPCServerStartingFailed.Error())
	}

	calcServer := NewServer(cfg, logger)

	lis, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%s",
		cfg.Microservice.Host,
		cfg.Microservice.Port,
	))
	if err != nil {
		logger.Fatal(neterr.ErrGRPCServerListeningFailed.Error())
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCalculatorServer(grpcServer, calcServer)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		grpcServer.GracefulStop()
	}()

	logger.Info("gRPC microservice starting: " + cfg.Microservice.Host + ":" + cfg.Microservice.Port)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(neterr.ErrGRPCServerServingFailed.Error())
	}

	return nil
}
