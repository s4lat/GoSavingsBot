-- +goose Up
CREATE TABLE expenses (
      id BIGSERIAL PRIMARY KEY,
      date TIMESTAMP NOT NULL,
      title VARCHAR(64) NOT NULL,
      amount DOUBLE PRECISION NOT NULL,
      currency VARCHAR(16) NOT NULL,
      type_id BIGINT DEFAULT NULL,
      user_id BIGINT NOT NULL
);

CREATE TABLE expense_types (
       id BIGSERIAL PRIMARY KEY,
       type_name VARCHAR(16) NOT NULL,
       type_color INTEGER DEFAULT 0,
       user_id BIGINT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS expense_types;
