create table if not exists notifications (
    id bigserial primary key,
    user_id text not null,
    pet_id uuid references pets (id) on delete set null,
    type text not null,
    detail jsonb not null default '{}'::jsonb,
    created_at timestamp with time zone default now(),
    seen_at timestamp with time zone default null
);

create index if not exists idx_notifications_user_id_created_at on notifications (user_id, created_at desc);
create index if not exists idx_notifications_spotted_pet on notifications ((detail->>'pet_id'), created_at)
    where type = 'spotted_pet';
