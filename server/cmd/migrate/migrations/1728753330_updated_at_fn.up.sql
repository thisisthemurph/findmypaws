create extension if not exists "uuid-ossp";

create or replace function fn_update_updated_at_timestamp()
    returns trigger as $$
begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language plpgsql;
