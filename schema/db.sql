create schema if not exists ng_lu;
use ng_lu;

create table if not exists accounts
(
    `username`        varchar(255) not null
        primary key,
    `hashed_password` varchar(255) not null default '',
    `google_id`       varchar(255) not null default '',
    `last_logout`     timestamp    not null default 0,
    `created_at`      timestamp             default CURRENT_TIMESTAMP null,
    `updated_at`      timestamp             default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
);

create table if not exists profiles
(
    `id`         int auto_increment
        primary key,
    `full_name`  varchar(255) not null default '',
    `phone`      varchar(255) not null default '',
    `email`      varchar(255) not null default '',
    `account_id` varchar(255) null,
    `created_at` timestamp             default CURRENT_TIMESTAMP null,
    `updated_at` timestamp             default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    index idx_account_id (account_id)
);
