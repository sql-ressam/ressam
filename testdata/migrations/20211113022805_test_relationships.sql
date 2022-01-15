-- +goose Up
-- +goose StatementBegin
create table books
(
    id   bigserial primary key,
    name varchar(255) not null
);

create table authors
(
    id   bigserial primary key,
    name varchar(255) not null unique
);

create table books_authors
(
    author_id bigint not null references authors,
    book_id   bigint not null references books,
    constraint books_authors_pk primary key (author_id, book_id)
);

create table users
(
    id   bigserial primary key,
    name varchar(255) null
);

create table books_users
(
    user_id        bigint    not null references users (id),
    book_id        bigint    not null references books (id),
    receiving_date timestamp not null default now(),
    constraint books_users_pk primary key (user_id, book_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists books, authors, books_authors, users, books_users cascade;
-- +goose StatementEnd
