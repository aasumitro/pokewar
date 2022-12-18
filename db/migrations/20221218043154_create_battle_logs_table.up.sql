CREATE TABLE IF NOT EXISTS battle_logs
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    battle_id INTEGER NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at INT,
    FOREIGN KEY(battle_id) REFERENCES battles(id)
)