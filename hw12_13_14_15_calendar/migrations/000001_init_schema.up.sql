create table if not exists event (
    id uuid primary key not null,
    title varchar not null,
    started_at timestamptz not null,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    description text
)