package task_manager

import (
	"sync"
	"context"
	"fmt"
	"time"

	pb "github.com/dinklen/GolangCalc_V2/api/proto/generated"
	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/google/uuid"
)

type TaskManager struct {
	client pb.CalculatorClient
	stream pb.Calculator_CalculateClient
	pending sync.Map
	mu sync.Mutex
	currentID int
	cfg *config.Config
}

func NewTaskManager(client *pb.CalculatorClient, cfg *config.Config) (*TaskManager, error) {
	stream, err := client.Calculate(context.Background())
	if err != nil {
		return nil, err
	}

	tm := &TaskManager{
		client: client,
		stream: stream,
		cfg: cfg,
	}

	go tm.receiver()

	return tm, nil
}

func (tm *TaskManager) receiver() {
	for {
		result, err := tm.stream.Recv()
		if err != nil {
			tm.pending.Range(func(key, value interface{}) bool {
				close(value.(chan *pb.Result))
				tm.pending.Delete(key)
				return true
			})
			return
		}

		if ch, ok := tm.pending.Load(result.Id); ok {
			ch.(chan *pb.Result) <- result
			tm.pending.Delete(result.Id)
		}
	}
}

func (tm *TaskManager) getNextID(parentID string) string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.currentID++
	return fmt.Sprintf("%s-%s", parentID, currentID)
}

func (tm *TaskManager) AsyncEvaluate(
	leftValue, rightValue float64,
	operator string,
	parentID string,
) (string, chan *pb.Result) {
	id := tm.getNexrID(parentID)
	ch := make(chan *pb.Result, 1)

	tm.pending.Store(id, ch)

	err := tm.stream.Send(&pb.Subexpression{
		Id: id,
		ParentId: parentId,
		LeftValue: leftValue,
		RightValue: rightValue,
		Operator: operator,
	})

	if err != nil {
		close(ch)
		return id, nil
	}
	return id, ch
}

func (tm *TaskManager) Evaluate(
	leftValue, rightValue float64,
	operator string,
	parentID string
) (float64, error) {
	id, ch := tm.AsyncEvaluate(leftValue, rightValue, operator, parentID)

	defer tm.currentID = 0

	select {
	case result := <-ch:
		if result.Error != nil {
			return 0.0, fmt.Errorf(result.Error)
		}
		return result.Value, nil
	case time.After(cfg.Microservice.WaitingTime):
		tm.pending.Delete(id)
		return 0.0, fmt.Errorf("timeout")
	}
}
