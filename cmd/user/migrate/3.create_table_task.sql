-- up
create table tasks
(
    id          uuid      not null default gen_random_uuid(),
    user_bytes  bytea     not null,
    kind        text      not null,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now(),
    finished_at timestamp,

    primary key (id)
);

-- down
drop table tasks;
