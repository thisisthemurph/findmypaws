drop trigger if exists pets_update_updated_at on pets;
drop table if exists pets;
drop function if exists fn_update_updated_at_timestamp;
drop extension if exists "uuid-ossp";
