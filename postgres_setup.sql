-- Create PostgreSQL database used for yl-pgwriter

\prompt 'Password for yolink (read/write) user: ' yolink_pw
\prompt 'Password for yolinkro (read only) user: ' yolinkro_pw

create database yolink;
\c yolink

create user yolink with encrypted password :'yolink_pw';
grant all privileges on database yolink to yolink;

create user yolinkro with encrypted password :'yolinkro_pw';
grant connect on database yolink to yolinkro;

set role yolink;

create table mqtt_messages(
    id bigserial,
    home_id text,
    dev_id text,
    time timestamp with time zone,
    data jsonb
);
create index dev_time on mqtt_messages(dev_id, time);

grant select on mqtt_messages to yolinkro;
