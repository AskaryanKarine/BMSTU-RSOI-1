-- +goose Up
-- +goose StatementBegin
create table if not exists persons (
    "id" serial primary key,
    "name" text not null,
    "age" int not null,
    "address" text not null,
    "work" text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists persons;
-- +goose StatementEnd
