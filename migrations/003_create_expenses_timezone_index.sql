-- +goose Up
CREATE INDEX expense_user_id_date_idx ON expenses (user_id, date)

-- +goose Down
DROP INDEX IF EXISTS expense_user_id_date_idx;
