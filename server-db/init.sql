\c nf6;

create table global_config (
  id integer not null primary key default 69,
  domain text not null,
  wireguard_public_key text not null,
  
  constraint singleton check (id = 69)
);

insert into global_config (domain, wireguard_public_key) values ('nf6.sh', '');

create table account (
  id bigserial primary key,
  email text not null unique,
  ssh_public_key text not null unique,
  ssl_public_key text not null unique
);

create table repo (
  id bigserial primary key,
  account_id bigint references account(id),
  name text not null,

  unique (account_id, name)
);

create table machine (
  id bigserial primary key,
  account_id bigint references account(id),
  host_name text not null,
  wireguard_public_key text not null unique,
  address_ipv6 inet not null unique,

  unique (account_id, host_name)
);
