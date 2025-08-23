-- drop table users;
-- drop table cryptos;
-- drop table prices;

create table users (
    username varchar(150) primary key not null,
    password varchar(200) not null
);


create table cryptos (
    id serial primary key,
    symbol varchar(50) not null,
    name varchar(50) not null,
    unique(symbol)
);

create table prices (
    id serial primary key,
    crypto int references cryptos (id) on delete cascade,
    price double precision not null,  
    updated_at timestamp not null
);


