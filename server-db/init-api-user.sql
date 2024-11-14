create user nf6_api with password 'PG_NF6_API_PASS';

grant usage on schema public to nf6_api;
grant usage on all sequences in schema public to nf6_api;
grant select, insert, update, delete on all tables in schema public to nf6_api;
