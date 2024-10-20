create table if not exists profiles (
    id uuid references auth.users (id) on delete cascade,
    name text not null,
    updated_at timestamp with time zone default now(),

    primary key (id)
);

create trigger profiles_update_updated_at
    before update on profiles
    for each row
execute function fn_update_updated_at_timestamp();
