CREATE TABLE IF NOT EXISTS battle_players
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    battle_id INTEGER NOT NULL,
    monster_id INTEGER NOT NULL,
    eliminated_at INT,
    annulled_at INT,
    rank INT NOT NULL,
    point INT NOT NULL,
    FOREIGN KEY(battle_id) REFERENCES battles(id),
    FOREIGN KEY(monster_id) REFERENCES monsters(id)
)