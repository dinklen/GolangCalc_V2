package handlers

import (
	"net/http"
	"sync"

	"github.com/dinklen/GolangCalc_V2/internal/service/calculator/parser"
	"github.com/dinklen/GolangCalc_V2/internal/seevice/calculator/evaluator"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/calc_errors"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

type Expression struct {
	Expr string `json:"expression"`
}

func CreateCalculatorHandler(db *sql.DB, 
func CalculatorHandler(ctx echo.Context) error {
	expression := new(Expression)
	if err := ctx.Bind(expression); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error":"invalid json format"})
	}

	if expression.Expr == "" {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error":"empty expression"})
	}

	tokens, err := parser.Tokenize(expression.Expr)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error":err.What()})
	}

	ast, err := parser.BuildAST(tokens)
	if err != nil {
                return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{"error":err.What()})
        }

	taskID := uuid.New().String()

	go func(){
		result, err := evaluator.Evaluate(ast, taskID)
	}()

	return ctx.JSON(http.StatusCreated, map[string]string{"taskID":taskID})
}
