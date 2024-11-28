package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/s4lat/gosavingsbot/internal/domain"
	"github.com/s4lat/gosavingsbot/internal/usecase"
	"net/http"
	"time"
)

type ExpenseRoutes struct {
	uc *usecase.UseCase
}

func NewExpenseRoutes(uc *usecase.UseCase) *ExpenseRoutes {
	return &ExpenseRoutes{
		uc: uc,
	}
}

func (r *ExpenseRoutes) getExpensesByDateAndTimeZone(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = getLoggerFromEchoContext(c)
		user = getServiceUser(c)
	)

	dateParam := c.QueryParam("date")

	date, err := time.Parse("2006-01-02", dateParam)
	if dateParam == "" || err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("date is invalid"))
	}

	expenses, err := r.uc.Expense.GetExpensesByDateAndTimeZone(
		ctx,
		user.Id,
		date,
	)
	if err != nil {
		log.Errorf("error on get expenses by date and tz: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, internalError)
	}

	expensesForResp := make([]expenseForResponse, 0, len(expenses))
	for _, expense := range expenses {
		expensesForResp = append(expensesForResp, expenseForResponse{
			Id:       expense.Id,
			Date:     expense.Date.Format("2006-01-02"),
			Title:    expense.Title,
			Amount:   expense.Amount,
			Currency: expense.Currency,
			TypeID:   expense.TypeID,
		})
	}

	return echo.NewHTTPError(http.StatusOK, getExpensesByDateAndTimeZoneResponse{Expenses: expensesForResp})
}

type getExpensesByDateAndTimeZoneResponse struct {
	Expenses []expenseForResponse `json:"expenses"`
}

type expenseForResponse struct {
	Id       int64   `json:"id"`
	Date     string  `json:"date"`
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	TypeID   int32   `json:"type_id"`
}

func (r *ExpenseRoutes) createExpense(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		log  = getLoggerFromEchoContext(c)
		user = getServiceUser(c)
		req  createExpenseRequest
	)

	if err := c.Bind(&req); err != nil {
		log.Warnf("error on binding: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("bad request"))
	}

	if req.Title == "" || len(req.Title) > 64 {
		return echo.NewHTTPError(http.StatusBadRequest,
			newErrorResponse("title must be from 1 to 64 characters long"))
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if req.Date == "" || err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("date invalid"))
	}

	if req.Currency == "" {
		return echo.NewHTTPError(http.StatusBadRequest, newErrorResponse("currency must be set"))
	}

	if _, ok := domain.GetCurrency(req.Currency); !ok {
		return echo.NewHTTPError(http.StatusBadRequest,
			newErrorResponseF("can't find currency '%s'", req.Currency))
	}

	expense, err := r.uc.Expense.CreateExpense(ctx, domain.Expense{
		Date:     date,
		Title:    req.Title,
		Amount:   req.Amount,
		Currency: req.Currency,
		TypeID:   0,
		UserId:   user.Id,
	})
	if err != nil {
		log.Errorf("error on create expense: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, internalError)
	}

	return echo.NewHTTPError(http.StatusOK, createExpenseResponse{
		Id:       expense.Id,
		Title:    expense.Title,
		Amount:   expense.Amount,
		Currency: expense.Currency,
		TypeID:   expense.TypeID,
		Date:     expense.Date.Format("2006-01-02"),
	})
}

type createExpenseRequest struct {
	Date     string  `json:"date" binding:"required"`
	Title    string  `json:"title" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Currency string  `json:"currency" binding:"required"`
}

type createExpenseResponse struct {
	Id       int64   `json:"id"`
	Title    string  `json:"title"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	TypeID   int32   `json:"type_id"`
	Date     string  `json:"date"`
}
