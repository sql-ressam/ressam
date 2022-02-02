-- +goose Up
-- +goose StatementBegin
create type test_enum as enum ('first', 'second');

create table test_default_values
(
    id                               bigserial primary key,

    int_null                         int       null,
    int_not_null                     int       not null,
    int_null_default_1               int       null     default 1,
    int_not_null_default_1           int       not null default 1,

    test_enum_null                   test_enum null,
    test_enum_not_null               test_enum not null,
    test_enum_null_default_first     test_enum null     default 'first',
    test_enum_not_null_default_first test_enum not null default 'first'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table test_default_values;
drop type test_enum;
-- +goose StatementEnd
