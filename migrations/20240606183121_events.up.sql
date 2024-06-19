CREATE TABLE if not exists events (

                        event_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                        id int PRIMARY KEY
);

insert into events (event_timestamp, id) values (current_timestamp,1);