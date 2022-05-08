alter table event add column if not exists user_id int not null;
alter table event add column if not exists finished_at timestamptz not null;
alter table event add column if not exists notify_interval bigint not null;