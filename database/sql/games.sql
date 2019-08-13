---
CREATE TABLE games (
    id varchar(255) primary key,
    song varchar(255),
    home_id varchar(255),
    home_score int,
    home_ready boolean,
    away_id varchar(255),
    away_score int,
    away_ready boolean,
    started timestamp without time zone,
    finished timestamp without time zone
);