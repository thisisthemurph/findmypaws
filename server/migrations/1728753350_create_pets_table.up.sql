create extension if not exists "uuid-ossp";

create or replace function fn_update_updated_at_timestamp()
    returns trigger as $$
begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language plpgsql;

create table if not exists pets (
    id uuid primary key default uuid_generate_v4(),
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
