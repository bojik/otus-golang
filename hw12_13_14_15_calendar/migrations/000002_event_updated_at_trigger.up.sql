create or replace function update_event_update_at()
    returns trigger AS
$$
begin
    NEW.updated_at = now();
    return NEW;
end ;
$$
language plpgsql;

create trigger update_updated_at BEFORE update on event FOR EACH ROW execute procedure  update_event_update_at();