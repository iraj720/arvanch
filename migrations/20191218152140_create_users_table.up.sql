create table if not exists users
(
    id              uuid        PRIMARY KEY,
    name            TEXT        not null CHECK (name <> ''),
    account_id      uuid        not null,
    created_at      timestamp   not null default now(),
    constraint fk_accounts
        foreign key(account_id)
            references accounts(id)
);

create index if not exists accounts_users_idx on accounts(id);