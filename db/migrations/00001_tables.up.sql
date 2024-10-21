create extension if not exists "uuid-ossp";

create table if not exists users
(
    id            uuid not null primary key default uuid_generate_v4(),
    login         text not null check (login = lower(login)),
    password_hash text not null,
    balance       int  not null             default 0 check (balance >= 0),
    withdrawn     int  not null             default 0 check (withdrawn >= 0)
);

create unique index if not exists login_unique on users (lower(login));

do $$ begin
    create type status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
exception
    when duplicate_object then null;
end $$;

create table if not exists orders
(
    id            text                        not null primary key,
    user_id       uuid                        not null references users (id),
    status        status                      not null default 'NEW',
    accrual       int                         not null default 0 check (accrual >= 0),
    uploaded_at   timestamp(0) with time zone not null,
    in_processing bool                        not null default false
);

create table if not exists withdrawals
(
    user_id      uuid                        not null references users (id),
    order_id     text                        not null primary key,
    sum          int                         not null check (sum > 0),
    processed_at timestamp(0) with time zone not null
);
