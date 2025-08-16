create table if not exists messages
(
    id          uuid        PRIMARY KEY,
    user_id     uuid        not null,
    payload     TEXT        not null CHECK (payload <> ''),
    language    VARCHAR(10) not null CHECK (language <> ''),
    created_at  timestamp   not null default now(),
    constraint fk_users
        foreign key(user_id)
            references users(id)
);

create index if not exists messages_users_idx on users(id);