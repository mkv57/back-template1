-- up
create table avatars
(
    id         uuid      not null,
    owner_id   uuid      not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    foreign key (owner_id) references users on delete cascade,
    unique (id),
    unique (owner_id, id)
);

alter table users
    add column current_avatar_id uuid;

-- down
alter table users
drop
column current_avatar_id;

drop table avatars;
