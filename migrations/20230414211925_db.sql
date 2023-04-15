-- +goose Up
CREATE TABLE change_currency (
                                 user_id SERIAL PRIMARY KEY ,
                                 chat_id BIGINT NOT NULL ,
                                 currency VARCHAR(3) NOT NULL,
                                 cost_currency DECIMAL(5,2) NOT NULL
);
CREATE TABLE currency (
                          user_id SERIAL PRIMARY KEY,
                          chat_id BIGINT NOT NULL,
                          currency VARCHAR(3) NOT NULL
);
-- +goose Down
DROP TABLE change_currency;
DROP TABLE currency;
