package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s4lat/gosavingsbot/internal/domain"
	"time"
)

type ExpenseRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseRepo(pool *pgxpool.Pool) *ExpenseRepo {
	return &ExpenseRepo{pool: pool}
}

// CreateExpense inserts a new expense and returns it.
func (r *ExpenseRepo) CreateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	query := `INSERT INTO expenses (date, title, amount, currency, type_id, user_id) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.pool.QueryRow(
		ctx,
		query,
		expense.Date,
		expense.Title,
		expense.Amount,
		expense.Currency,
		expense.TypeID,
		expense.UserId,
	).Scan(&expense.Id)
	if err != nil {
		return domain.Expense{}, err
	}
	return expense, nil
}

// UpdateExpense updates an existing expense and returns it.
func (r *ExpenseRepo) UpdateExpense(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	query := `UPDATE expenses SET date = $1, title = $2, amount = $3, currency = $4, type_id = $5 
			  WHERE id = $6 AND user_id = $7 RETURNING id, date, title, amount, currency, type_id, user_id`
	err := r.pool.QueryRow(
		ctx,
		query,
		expense.Date,
		expense.Title,
		expense.Amount,
		expense.Currency,
		expense.TypeID,
		expense.Id,
		expense.UserId,
	).Scan(
		&expense.Id,
		&expense.Date,
		&expense.Title,
		&expense.Amount,
		&expense.Currency,
		&expense.TypeID,
		&expense.UserId,
	)
	if err != nil {
		return domain.Expense{}, err
	}
	return expense, nil
}

// DeleteExpense deletes an expense by ID.
func (r *ExpenseRepo) DeleteExpense(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// CreateExpenseType inserts a new expense type and returns it.
func (r *ExpenseRepo) CreateExpenseType(
	ctx context.Context,
	expenseType domain.ExpenseType,
) (domain.ExpenseType, error) {
	query := `INSERT INTO expense_types (type_name, type_color, user_id) 
			  VALUES ($1, $2, $3) RETURNING id`
	err := r.pool.QueryRow(
		ctx,
		query,
		expenseType.TypeName,
		expenseType.TypeColor,
		expenseType.UserId,
	).Scan(&expenseType.Id)
	if err != nil {
		return domain.ExpenseType{}, err
	}
	return expenseType, nil
}

// UpdateExpenseType updates an expense type after checking existence and returns it.
func (r *ExpenseRepo) UpdateExpenseType(
	ctx context.Context,
	expenseType domain.ExpenseType,
) (domain.ExpenseType, error) {
	checkQuery := `SELECT COUNT(1) FROM expense_types WHERE id = $1 AND user_id = $2`
	var count int
	err := r.pool.QueryRow(ctx, checkQuery, expenseType.Id, expenseType.UserId).Scan(&count)
	if err != nil || count == 0 {
		return domain.ExpenseType{}, errors.New("expense type does not exist")
	}

	updateQuery := `UPDATE expense_types SET type_name = $1, type_color = $2 
					WHERE id = $3 AND user_id = $4 RETURNING id, type_name, type_color, user_id`
	err = r.pool.QueryRow(
		ctx,
		updateQuery,
		expenseType.TypeName,
		expenseType.TypeColor,
		expenseType.Id,
		expenseType.UserId,
	).Scan(&expenseType.Id, &expenseType.TypeName, &expenseType.TypeColor, &expenseType.UserId)
	if err != nil {
		return domain.ExpenseType{}, err
	}
	return expenseType, nil
}

// DeleteExpenseType deletes an expense type by ID for a specific user.
func (r *ExpenseRepo) DeleteExpenseType(ctx context.Context, id int64, userId int64) error {
	query := `DELETE FROM expense_types WHERE id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, query, id, userId)
	return err
}

func (r *ExpenseRepo) GetExpensesCountByDateWithTz(ctx context.Context, userId int64, date time.Time) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ExpenseRepo) GetExpensesByDateWithTz(
	ctx context.Context,
	userId int64,
	date time.Time,
) ([]domain.Expense, error) {
	query := `SELECT 
					id,
					date,
					title, 
					amount, 
					currency, 
					type_id, 
					user_id
				FROM expenses
				WHERE user_id = $1 AND date = $2
    			ORDER BY created_at;`

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	rows, err := r.pool.Query(
		ctx,
		query,
		userId,
		date.Format("2006-01-02"),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]domain.Expense, 0), nil
		}
		return nil, wrapError(err)
	}

	expenses := make([]domain.Expense, 0, 100)
	for rows.Next() {
		var expense domain.Expense
		if err := rows.Scan(
			&expense.Id,
			&expense.Date,
			&expense.Title,
			&expense.Amount,
			&expense.Currency,
			&expense.TypeID,
			&expense.UserId,
		); err != nil {
			return nil, wrapError(err)
		}

		//expense.Date = expense.Date.In(tzLoc)
		expenses = append(expenses, expense)
	}

	return expenses, nil
}
