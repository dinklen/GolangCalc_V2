package database

import (
	"database/sql"
	"fmt"

	"github.com/dinklen/GolangCalc_V2/internal/models"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/dberr"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CreateExpression(db *sql.DB, expression string, userID string, logger *zap.Logger) (string, error) {
	var expressionID uuid.UUID

	err := db.QueryRow(
		`
		INSERT INTO expressions(user_id, expression, status)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		userID, expression, "pending",
	).Scan(&expressionID)

	if err != nil {
		return "", dberr.ErrExpressionCreatingFailed
	}

	logger.Info("expression creating success")
	return expressionID.String(), nil
}

func UpdateExpressionState(db *sql.DB, expressionID string, status string, value float64, logger *zap.Logger) error {
	_, err := db.Exec(
		`
		UPDATE expressions
		SET
			status = $1,
			result = $2
		WHERE id = $3
		`,
		status, value, expressionID,
	)

	if err != nil {
		return dberr.ErrExpressionUpdatingStateFailed
	}

	logger.Info("expression state updating success")
	return nil
}

func GetExpressions(db *sql.DB, userID string, logger *zap.Logger) ([]models.Expression, error) {
	rows, err := db.Query(
		`
		SELECT id, expression, result, status, created_at FROM expressions
		WHERE user_id = $1
		`,
		userID,
	)

	if err != nil {
		return []models.Expression{}, dberr.ErrExpressionsGettingFailed
	}
	logger.Info("expressions select success")

	defer rows.Close()

	var expressions []models.Expression
	for rows.Next() {
		var expression models.Expression
		if err := rows.Scan(
			&expression.ID,
			&expression.Expr,
			&expression.Result,
			&expression.Status,
			&expression.CreatingTime,
		); err != nil {
			return []models.Expression{}, dberr.ErrExpressionsParsingFailed
		}
		expressions = append(expressions, expression)
	}

	if err = rows.Err(); err != nil {
		return []models.Expression{}, dberr.ErrExpressionsParsingFailed
	}
	logger.Info("expressions getting success")

	return expressions, nil
}

func GetCurrentExpression(db *sql.DB, userID, expressionID string, logger *zap.Logger) (models.Expression, error) {
	var expression models.Expression

	err := db.QueryRow(
		`
		SELECT id, expression, result, status, created_at FROM expressions
		WHERE user_id = $1 AND id = $2
		`,
		userID, expressionID,
	).Scan(
		&expression.ID,
		&expression.Expr,
		&expression.Result,
		&expression.Status,
		&expression.CreatingTime,
	)

	if err != nil {
		return models.Expression{}, dberr.ErrExpressionGettingFailed
	}

	logger.Info("expression getting success")
	return expression, nil
}

func CreateSubexpression(db *sql.DB, parentExpressionID, subexpressionID string, subexpression string, logger *zap.Logger) (string, error) {
	var newSubexpressionID string

	parentExpressionUUID, err := uuid.Parse(parentExpressionID)
	if err != nil {
		return "", dberr.ErrParentExpressionIDParsingFailed
	}

	newID := fmt.Sprintf("%v-%v", parentExpressionID, subexpressionID)

	err = db.QueryRow(
		`
		INSERT INTO sub_expressions (id, expression_id, sub_expression, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`,
		newID, parentExpressionUUID, subexpression, "pending",
	).Scan(&newSubexpressionID)

	if err != nil {
		fmt.Println(err.Error())
		return "", dberr.ErrSubexpressionCreatingFailed
	}
	logger.Info("subexpression creating success")

	return newSubexpressionID, nil
}

func UpdateSubexpressionState(db *sql.DB, subexpressionID string, status string, value float64, logger *zap.Logger) error {
	_, err := db.Exec(
		`
		UPDATE sub_expressions
		SET
			status = $1,
			result = $2
		WHERE id = $3
		`,
		status, value, subexpressionID,
	)

	if err != nil {
		return dberr.ErrSubexpressionUpdatingStateFailed
	}

	logger.Info("subexpression updating success")
	return nil
}
