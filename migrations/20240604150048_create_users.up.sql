create table IF NOT EXISTS users(
    id bigserial not null primary key,
    email varchar not null unique ,
    passwordHash varchar not null,
    nme varchar not null,
    birthday date not null,
    calendar_ids INTEGER[]

);