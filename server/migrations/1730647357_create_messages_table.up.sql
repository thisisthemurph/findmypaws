create table if not exists conversations (
    id bigserial primary key,
    identifier uuid references pets (id) on delete cascade,
    primary_participant_id text not null,
    secondary_participant_id text not null,
    last_message_at timestamp with time zone,
    created_at timestamp with time zone not null default now(),
    unique (identifier, secondary_participant_id)
);

create index if not exists idx_conversations_identifier on conversations (identifier);
create index if not exists idx_conversations_primary_participant_id on conversations (primary_participant_id);
create index if not exists idx_conversations_secondary_participant_id on conversations (secondary_participant_id);

create table if not exists messages (
    id bigserial primary key,
    conversation_id bigserial references conversations (id) on delete cascade,
    sender_id text not null,
    text text not null check (char_length(text) <= 500),
    created_at timestamp with time zone not null default now(),
    read_at timestamp with time zone default null
);

create index idx_messages_conversation_id on messages (conversation_id, created_at);
create index idx_messages_sender_id on messages (sender_id);

create or replace function fn_conversations_update_last_message_at()
    returns trigger as $$
begin
    update conversations
    set last_message_at = now()
    where id = new.conversation_id;

    return new;
end;
$$ language plpgsql;

create trigger tr_messages_set_last_message_at
    after insert on messages
    for each row
execute function fn_conversations_update_last_message_at();
