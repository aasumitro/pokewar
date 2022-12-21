CREATE TABLE IF NOT EXISTS monsters
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    origin_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    base_exp INT NOT NULL,
    height INT NOT NULL,
    weight INT NOT NULL,
    avatar VARCHAR(255) NOT NULL,
    types TEXT NOT NULL,
    stats TEXT NOT NULL,
    skills TEXT NOT NULL,
    created_at INT,
    updated_at INT
)