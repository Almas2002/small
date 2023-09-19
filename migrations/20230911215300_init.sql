-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products
(
    product_id serial primary key,
    title      varchar(255) unique not null,
    price      float8              not null
);

CREATE TABLE IF NOT EXISTS users
(
    user_id serial primary key,
    phone   varchar(20),
    email   varchar(100)
);

CREATE TABLE IF NOT EXISTS user_sub_products
(
    user_id    int,
    product_id int,
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (user_id),
    CONSTRAINT product_fk FOREIGN KEY (product_id) REFERENCES products (product_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sub_products;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
