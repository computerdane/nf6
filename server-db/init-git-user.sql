create user nf6_git with password 'PG_NF6_GIT_PASS';

grant usage on schema public to nf6_git;
grant usage on all sequences in schema public to nf6_git;
grant select on table account to nf6_git;
grant select on table repo to nf6_git;
