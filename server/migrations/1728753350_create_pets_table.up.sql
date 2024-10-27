create table if not exists pets (
    id uuid primary key default uuid_generate_v4(),
    user_id text,
    type varchar(16),
    name text not null,
    tags jsonb,
    dob date,
    avatar_uri text,
    blurb text,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create index idx_pets_user_id on pets (user_id);

create trigger pets_update_updated_at
    before update on pets
    for each row
execute function fn_update_updated_at_timestamp();
