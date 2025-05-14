package handlers

import (
	"database/sql"
	"net/http"

	"github.com/dinklen/GolangCalc_V2/internal/database"
	"github.com/dinklen/GolangCalc_V2/internal/models"
	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/evaluator"
	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/tm"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func CreateCalculatorHandler(db *sql.DB, taskManager *tm.TaskManager, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		expression := new(models.Expression)
		if err := ctx.Bind(expression); err != nil {
			logger.Error("binding failed")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json format"})
		}

		if expression.Expr == "" {
			logger.Error("empty expression got")
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "empty expression"})
		}

		tokens, err := parser.Tokenize(expression.Expr, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		}

		ast, err := parser.BuildAST(tokens, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		}

		userID, ok := ctx.Get("userID").(string)
		if !ok {
			logger.Error("failed to get user ID")
			return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found")
		}

		taskID, dberr := database.CreateExpression(db, expression.Expr, userID, logger)
		if dberr != nil {
			logger.Error(dberr.Error())
			ctx.JSON(http.StatusInternalServerError, map[string]string{"error": dberr.Error()})
		}

		go func() {
			result, err := evaluator.Evaluate(ast, taskManager, taskID, logger)
			if err != nil {
				err := database.UpdateExpressionState(db, taskID, "failed", 0.0, logger)
				if err != nil {
					logger.Error(err.Error())
				}
			} else {
				err := database.UpdateExpressionState(db, taskID, "calculated", result, logger)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		}()

		logger.Info("success to creating task")
		return ctx.JSON(http.StatusCreated, map[string]string{"task_id": taskID})
	}
}

func CreateExpressionsHandler(db *sql.DB, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, ok := ctx.Get("userID").(string)
		if !ok {
			logger.Error("failed to get user ID")
			return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found")
		}

		expressions, err := database.GetExpressions(db, userID, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		logger.Info("success to handle expression")
		return ctx.JSON(http.StatusOK, map[string][]models.Expression{"expressions": expressions})
	}
}

func CreateCurrentExpressionHandler(db *sql.DB, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		IDParam := ctx.Param("id")

		userID, ok := ctx.Get("userID").(string)
		if !ok {
			logger.Error("failed to get user ID")
			return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found")
		}

		expression, err := database.GetCurrentExpression(db, userID, IDParam, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get expression"})
		}

		logger.Info("success to handle current expression")
		return ctx.JSON(http.StatusOK, map[string]models.Expression{"expression": expression})
	}
}
