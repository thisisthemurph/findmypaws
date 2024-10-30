create table if not exists alerts (
    id bigserial primary key,
    pet_id uuid references pets (id) on delete cascade,
    user_id text, -- ID of the user creating the alert
    anonymous_user_id text, -- ID of the anonymous user creating the alert
    created_at timestamp with time zone default now()
);

create index if not exists idx_alerts_pet_id on alerts (pet_id);
create index if not exists idx_alerts_user_id on alerts (user_id);
create index if not exists idx_alerts_anonymous_user_id on alerts (anonymous_user_id);
