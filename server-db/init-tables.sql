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
  wg_public_key text not null unique check (wg_public_key <> ''),
  addr_ipv6 inet not null unique,

  unique (account_id, host_name)
);

create table subnet (
  id bigserial primary key,
  account_id bigint references account(id),
  subnet_ipv6 cidr not null unique
);
