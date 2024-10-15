create table if not exists pets (
    id uuid primary key default uuid_generate_v4(),
    user_id uuid references profiles (id) on delete cascade,
    type varchar(16),
    name text not null,
    tags jsonb,
    dob date,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create trigger pets_update_updated_at
    before update on pets
    for each row
execute function fn_update_updated_at_timestamp();
