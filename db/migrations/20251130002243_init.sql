-- +goose Up
-- +goose StatementBegin
CREATE TYPE account_group AS ENUM ('Bank', 'Investment', 'AssetOrLoan');
CREATE TYPE account_type AS ENUM ('BankingAccount', 'InvestmentAccount');
CREATE TYPE category_type AS ENUM ('Income', 'Expense', 'Transfer');
CREATE TYPE action_type AS ENUM ('cont', 'int', 'with', 'buy', 'sell', 'div', 'add', 'fee', 'drip');
CREATE TYPE security_type AS ENUM ('stock', 'bond', 'mutualfund', 'etf', 'crypto');

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT not null,
    account_group account_group not null,
    hidden boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null,
    user_uuid UUID not null,
    type account_type not null
);

CREATE TABLE IF NOT EXISTS categories(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT not null,
    category_type category_type not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null,
    user_uuid UUID not null,
    description TEXT null,
    hidden boolean not null default false
);

CREATE TABLE IF NOT EXISTS banking_transactions(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    date timestamp not null,
    amount decimal(10, 4) not null,
    payee TEXT null,
    memo TEXT null,
    account_id UUID not null REFERENCES accounts(id),
    category_id UUID not null REFERENCES categories(id),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null,
    balance decimal(10, 4) not null
);

CREATE TABLE IF NOT EXISTS currency_conversions(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    date timestamp not null,
    pair TEXT not null,
    amount decimal(10, 4) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null
);

CREATE TABLE IF NOT EXISTS securities(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT not null,
    symbol TEXT not null,
    security_type security_type not null,
    currency TEXT not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null
);

CREATE TABLE IF NOT EXISTS investment_transactions(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    date timestamp not null,
    action action_type not null,
    security_id UUID not null REFERENCES securities(id),
    price decimal(10, 4) not null,
    memo TEXT null,
    category_id UUID not null REFERENCES categories(id),
    amount decimal(10, 4) not null,
    account_id UUID not null REFERENCES accounts(id),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null,
    cash_balance decimal(10, 4) not null,
    share_balance decimal(10, 4) not null
);

CREATE TABLE IF NOT EXISTS prices(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    date timestamp not null,
    security_id UUID not null REFERENCES securities(id),
    price decimal(10, 4) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp null
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS prices;
DROP TABLE IF EXISTS investment_transactions;
DROP TABLE IF EXISTS securities;
DROP TABLE IF EXISTS currency_conversions;
DROP TABLE IF EXISTS banking_transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;

DROP TYPE IF EXISTS security_type;
DROP TYPE IF EXISTS action_type;
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS account_type;
DROP TYPE IF EXISTS account_group;
-- +goose StatementEnd
