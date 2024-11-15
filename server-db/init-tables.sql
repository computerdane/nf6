create table global_config (
  id integer not null primary key default 69,
  domain text not null check (domain <> ''),
  wireguard_public_key text not null check (wireguard_public_key <> ''),
  
  constraint singleton check (id = 69)
);

insert into global_config (domain, wireguard_public_key) values ('nf6.sh', 'sample');

create table account (
  id bigserial primary key,
  email text not null unique check (email <> ''),
  ssh_public_key text not null unique check (ssh_public_key <> ''),
  ssl_public_key text not null unique check (ssl_public_key <> '')
);

create table repo (
  id bigserial primary key,
  account_id bigint references account(id),
  name text not null check (name <> '') check (name ~* '^[A-Za-z0-9][A-Za-z0-9\-_]+[A-Za-z0-9]$'),

  unique (account_id, name)
);

create table machine (
  id bigserial primary key,
  account_id bigint references account(id),
  host_name text not null check (host_name <> ''),
  wireguard_public_key text not null unique check (wireguard_public_key <> ''),
  address_ipv6 inet not null unique,

  unique (account_id, host_name)
);
