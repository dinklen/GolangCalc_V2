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
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCalculatorServer
	workerPool *WorkerPool
	config     *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config:     cfg,
		workerPool: NewWorkerPool(cfg.Microservice.ComputingPower),
	}
}

type WorkerPool struct {
	tasks    chan *pb.Subexpression
	results  chan *pb.Result
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewWorkerPool(quantity int) *WorkerPool {
	pool := &WorkerPool{
		tasks:    make(chan *pb.Subexpression),
		results:  make(chan *pb.Result),
		shutdown: make(chan struct{}),
	}

	pool.wg.Add(quantity)
	for i := 0; i < quantity; i++ {
		go pool.worker()
	}

	return pool
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()

	for {
		select {
		case task := <-p.tasks:
			res, err := evaluateExpression(task.RightValue, task.LeftValue, task.Operator)
			response := &pb.Result{Id: task.Id}
			if err != nil {
				response.Error = err.Error()
			} else {
				response.Value = res
			}
			p.results <- response

		case <-p.shutdown:
			return
		}
	}
}

func (p *WorkerPool) Shutdown() {
	close(p.shutdown)
	p.wg.Wait()
	close(p.results)
}

func evaluateExpression(leftValue, rightValue float64, operator string) (float64, error) {
	result, err := evaluator.Calculate(leftValue, rightValue, operator)
	if err != nil {
		return 0.0, fmt.Errorf("%s", err.What())
	}

	return result, nil
}

func (s *Server) Calculate(stream pb.Calculator_CalculateServer) error {
	defer s.workerPool.Shutdown()

	go func() {
		for result := range s.workerPool.results {
			if err := stream.Send(result); err != nil {
				fmt.Printf("error to send response: %v", err)
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Printf("client connection closed")
			return nil
		}

		select {
		case s.workerPool.tasks <- req:
		default:
			s.workerPool.results <- &pb.Result{
				Id:    req.Id,
				Error: "server is busy",
			}
		}
	}
}

func Start(cfg *config.Config) {
	calcServer := NewServer(cfg)

	lis, err := net.Listen("tcp", cfg.Microservice.Port)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, calcServer)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		grpcServer.GracefulStop()
	}()

	fmt.Printf("Server started on :%s", cfg.Microservice.Port)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
