package evaluator

import (
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calcerr"

	"go.uber.org/zap"
)

func Calculate(leftVal, rightVal float64, operator string, cfg *config.Config, logger *zap.Logger) (float64, error) {
	switch operator {
	case "+":
		time.Sleep(cfg.Microservice.PlusTime)
		logger.Info("operation plus finished successful")
		return leftVal + rightVal, nil

	case "-":
		time.Sleep(cfg.Microservice.MinusTime)
		logger.Info("operation minus finished successful")
		return leftVal - rightVal, nil

	case "*":
		time.Sleep(cfg.Microservice.MultiplyTime)
		logger.Info("operation multiply finished successful")
		return leftVal * rightVal, nil

	case "/":
		time.Sleep(cfg.Microservice.DivisionTime)
		if rightVal == 0 {
			return 0, calcerr.ErrDivisionByZero
		}

		logger.Info("operation division finished successful")
		return leftVal / rightVal, nil

	case "^":
		time.Sleep(cfg.Microservice.PowerTime)
		result := 1.0
		for i := 0; i < int(rightVal); i++ {
			result *= leftVal
		}

		logger.Info("operation power finished successful")
		return result, nil

	default:
		return 0, calcerr.ErrUnknownOperator
	}
}
