-- +migrate Up

create table if not exists responses (
    id uuid primary key,
    status text not null,
    error text,
    payload jsonb,
    created_at timestamp without time zone not null default current_timestamp
);

create table if not exists users (
    id bigint unique,
    mail_id text primary key,
    email text not null,
    name text not null,
    updated_at timestamp with time zone not null default current_timestamp,
    created_at timestamp with time zone default current_timestamp
);

create index if not exists users_id_idx on users(id);
create index if not exists users_email_idx on users(email);
create index if not exists users_mailid_idx on users(mail_id);

create table if not exists links (
    id serial primary key,
    link text not null,
    unique(link)
);

create index if not exists links_link_idx on links(link);

create table if not exists permissions (
    request_id text not null,
    mail_id text not null,
    link text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null default current_timestamp,

    unique (mail_id, link),
    foreign key(mail_id) references users(mail_id) on delete cascade on update cascade,
    foreign key(link) references links(link) on delete cascade on update cascade
);

create index if not exists permissions_mailid_idx on permissions(mail_id);
create index if not exists permissions_link_idx on permissions(link);

-- +migrate Down

drop index if exists permissions_mailid_idx;
drop index if exists permissions_link_idx;

drop table if exists permissions;

drop index if exists links_link_idx;

drop table if exists links;

drop index if exists users_id_idx;
drop index if exists users_email_idx;
drop index if exists users_mailid_idx;

drop table if exists users;
drop table if exists responses;