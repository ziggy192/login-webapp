create schema if not exists ng_lu;
use ng_lu;

create table if not exists accounts
(
    username        varchar(255)                           not null
        primary key,
    hashed_password varchar(255) default ''                not null,
    google_id       varchar(255) default ''                not null,
    created_at      timestamp    default CURRENT_TIMESTAMP null,
    updated_at      timestamp    default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);


create table if not exists profiles
(
    id         int auto_increment
        primary key,
    full_name  varchar(255) default ''                not null,
    phone      varchar(255) default ''                not null,
    email      varchar(255) default ''                not null,
    account_id varchar(255)                           null,
    created_at timestamp    default CURRENT_TIMESTAMP null,
    updated_at timestamp    default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint idx_uq_account_id
        unique (account_id)
);
