---
CREATE TABLE scores (
    game VARCHAR(255) DEFAULT '' ,
    player VARCHAR(255) DEFAULT '',
    points INT DEFAULT 0,
    PRIMARY KEY (game, player)
);