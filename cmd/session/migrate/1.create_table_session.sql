-- up
create table sessions
(
    id         uuid      not null,
    token      text      not null,
    ip         text      not null,
    user_agent text      not null,
    user_id    uuid      not null,
    status     text      not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    unique (token),
    primary key (id)
);

create type deduplication_kind as enum ('StatusUpdate');

create table deduplication
(
    id         uuid               not null,
    kind       deduplication_kind not null,
    created_at timestamp          not null default now(),
    primary key (id, kind)
);

-- down
drop table sessions;
drop type deduplication_kind;
drop table deduplication;
