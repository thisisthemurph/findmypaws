create table if not exists anonymous_users (
    id text not null primary key,
    name text,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create trigger anonymous_users_update_updated_at
    before update on anonymous_users
    for each row
execute function fn_update_updated_at_timestamp();
