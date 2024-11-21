create table if not exists account (
    id bigserial primary key,
    email text not null unique check (email <> ''),
    ssh_pub_key text not null unique check (ssh_pub_key <> ''),
    tls_pub_key text not null unique check (tls_pub_key <> '')
);

create table if not exists host (
    id bigserial primary key,
    account_id bigint references account(id),
    name text not null check (name ~* '^[A-Za-z0-9]([A-Za-z0-9_-]{0,61}[A-Za-z0-9])?$'),
    addr6 inet not null unique,
    wg_pub_key text not null unique check (wg_pub_key <> ''),
    tls_pub_key text not null unique check (tls_pub_key <> ''),

    unique (account_id, name)
);

create table if not exists repo (
    id bigserial primary key,
    account_id bigint references account(id),
    name text not null check (name <> '') check (name ~* '^[A-Za-z0-9]+[A-Za-z0-9\-_]+[A-Za-z0-9]+$'),

    unique (account_id, name)
);
