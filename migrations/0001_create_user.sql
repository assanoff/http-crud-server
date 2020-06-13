create schema if not exists {{.Schema}};
create sequence if not exists {{.Schema}}.seq_users;
create table if not exists {{.Schema}}.users
(
  id int not null default nextval('{{.Schema}}.seq_users'::regclass),
  name varchar not null,
  email varchar not null,
  constraint "PK_users" primary key (id),
  constraint "UQ_users_email" unique (email),
  constraint "CHK_users_email" check (email like '%@%')
);
 ---- create above / drop below ----
