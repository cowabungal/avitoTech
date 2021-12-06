CREATE TABLE users
(
    user_id integer not null unique,
    balance float   not null
);

CREATE TABLE transactions
(
    id serial not null unique,
    user_id integer not null,
    operation varchar(40) not null,
    date timestamp not null
);