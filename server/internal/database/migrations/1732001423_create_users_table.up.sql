create table if not exists users (
    id text primary key,
    data jsonb not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create trigger users_update_updated_at
    before update on users
    for each row
execute function fn_update_updated_at_timestamp();
