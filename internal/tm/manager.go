package tm

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	pb "github.com/dinklen/GolangCalc_V2/api/proto/generated"
	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/database"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/tmerr"

	"go.uber.org/zap"
)

type TaskManager struct {
	client    pb.CalculatorClient
	stream    pb.Calculator_CalculateClient
	pending   sync.Map
	mu        sync.Mutex
	currentID int
	cfg       *config.Config
	db        *sql.DB
}

func NewTaskManager(client pb.CalculatorClient, cfg *config.Config, db *sql.DB, logger *zap.Logger) (*TaskManager, error) {
	stream, err := client.Calculate(context.Background())
	if err != nil {
		return nil, tmerr.ErrStreamCreatingFailed
	}
	logger.Info("stream creating success")

	taskManager := &TaskManager{
		client: client,
		stream: stream,
		db:     db,
		cfg:    cfg,
	}

	go taskManager.receiver(logger)

	logger.Info("task manager creating success")
	return taskManager, nil
}

func (taskManager *TaskManager) receiver(logger *zap.Logger) {
	logger.Info("task receiving started")
	for {
		result, err := taskManager.stream.Recv()
		if err != nil {
			logger.Error(err.Error())
			taskManager.pending.Range(func(key, value interface{}) bool {
				close(value.(chan *pb.Result))
				taskManager.pending.Delete(key)
				return true
			})
			return
		}

		if ch, ok := taskManager.pending.Load(result.Id); ok {
			logger.Info("task result got")
			ch.(chan *pb.Result) <- result
			taskManager.pending.Delete(result.Id)
		}
	}
}

func (taskManager *TaskManager) getNextID(parentID string) string {
	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	taskManager.currentID++
	return fmt.Sprintf("%s-%d", parentID, taskManager.currentID)
}

func (taskManager *TaskManager) AsyncEvaluate(leftValue, rightValue float64, operator string, parentID string, logger *zap.Logger) (string, chan *pb.Result) {
	logger.Info("async evaluate started")
	id := taskManager.getNextID(parentID)
	ch := make(chan *pb.Result, 1)

	taskManager.pending.Store(id, ch)

	err := taskManager.stream.Send(&pb.Subexpression{
		Id:         id,
		ParentId:   parentID,
		LeftValue:  leftValue,
		RightValue: rightValue,
		Operator:   operator,
	})

	if err != nil {
		logger.Error(err.Error())
		close(ch)
		return id, nil
	}
	return id, ch
}

func (taskManager *TaskManager) Evaluate(leftValue, rightValue float64, operator string, parentID string, logger *zap.Logger) (float64, error) {
	logger.Info("task evaluating started")
	subexprID, err := database.CreateSubexpression(
		taskManager.db,
		parentID, fmt.Sprintf("%d", taskManager.currentID),
		fmt.Sprintf("%v%v%v", leftValue, operator, rightValue),
		logger,
	)

	if err != nil {
		return 0.0, err
	}

	id, ch := taskManager.AsyncEvaluate(leftValue, rightValue, operator, parentID, logger)

	select {
	case result := <-ch:
		logger.Info("result got")
		if result.Error != "" {
			err := database.UpdateSubexpressionState(taskManager.db, subexprID, "failed", 0.0, logger)
			if err != nil {
				return 0.0, err
			}
			return 0.0, fmt.Errorf(result.Error)
		}

		err := database.UpdateSubexpressionState(taskManager.db, subexprID, "calculated", result.Value, logger)
		if err != nil {
			return 0.0, err
		}
		return result.Value, nil

	case <-time.After(taskManager.cfg.Microservice.WaitingTime):
		err := database.UpdateSubexpressionState(taskManager.db, subexprID, "failed", 0.0, logger)
		if err != nil {
			return 0.0, err
		}

		taskManager.pending.Delete(id)
		return 0.0, tmerr.ErrTimeout
	}
}
