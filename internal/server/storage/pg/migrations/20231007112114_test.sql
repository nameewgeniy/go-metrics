-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics_gauge
(
    name  VARCHAR(255)     NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

CREATE TABLE metrics_counter
(
    name  VARCHAR(255) NOT NULL PRIMARY KEY,
    value BIGINT       NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics_gauge;
DROP TABLE IF EXISTS metrics_counter;
-- +goose StatementEnd
