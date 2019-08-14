---
CREATE TABLE games (
    id VARCHAR(255) PRIMARY KEY,
    song VARCHAR(255) DEFAULT 'Rick Astley - Never gonna give you up',
    home_id VARCHAR(255) DEFAULT '',
    home_score INT DEFAULT 0,
    home_ready BOOLEAN DEFAULT 'FAlse',
    away_id VARCHAR(255) DEFAULT '',
    away_score INT DEFAULT 0,
    away_ready BOOLEAN DEFAULT 'FAlse',
    started BIGINT DEFAULT 0,
    finished BIGINT DEFAULT 0
);