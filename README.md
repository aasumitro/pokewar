# Pokewar

Pocket Monster Battleroyale

<p align="center">
<img src="./assets/battleroyale.gif" alt="gif">
</p>

---
### NTK - Battleroyale

1. 5 players in 1 battle/round
2. Record battle history
3. Battle ranks & scores 
4. Monsters point based on battle (accumulate rank)
5. Annulled player rank (in the end of battle)

---
#### Basic Entity Info
1. Monsters:
    
   List of Available Monsters (Dex)
   (id, origin_id, name/species, base_exp, height, weight, avatar, types, stats, skills, created_at, updated_at)

2. Battles: 

    Battles Arena with Random/Selected(?) Monster
    (id, started_at, ended_at)

3. Players: 

    List of Played Monsters in 1 Battle 
    (id, battle_id, monster_id, eliminate_at, annulled_at, rank, point)

4. Logs:

    List of Battles Log 
    (id, battle_id, description, created_at)
   
...

#### Entity (*)
```mermaid
erDiagram
    MONSTERS {
        id int  
        origin_id int
        base_exp int
        height int
        weight int 
        avatar string
        types text
        stats text
        skills text 
        created_at int
        updated_at int
    }
    
    BATTLES {
        id int
        started_at int
        ended_at int
    }
    PLAYERS {
        id int
        battle_id int
        monster_id int
        eliminated_at int 
        annulled_at int
        rank int
        point int
    }
    LOGS {
        id int
        battle_id int
        description string
        created_at int
    }
    
    BATTLES }|--|{ PLAYERS : m2m
    BATTLES ||--|{ LOGS : o2m
    MONSTERS }|--|{ PLAYERS : m2m
```

#### Sequence (*)
```mermaid
   sequenceDiagram
        participant DELIVERY
        participant SERVICE
        participant REST_REPOSITORY
        participant SQL_REPOSITORY
        par Data Flow
            DELIVERY->>+SERVICE: ...
        and
            SERVICE->>+REST_REPOSITORY: ...
        and
            SERVICE->>+SQL_REPOSITORY: ...
        end 
   
     
        REST_REPOSITORY-->>-SERVICE: ...
        Note over REST_REPOSITORY,SERVICE: ...
   
        SQL_REPOSITORY-->>-SERVICE: ...
        Note over SQL_REPOSITORY,SERVICE: ...
        
        SERVICE-->>-DELIVERY: ...
        Note over SERVICE,DELIVERY: ...     
   ```

---
### GTK
Lorem Ipsum is simply dummy text of the printing and typesetting industry. 

#### Required
- Sqlite3
- Golang v1.19

#### Database Migration
###### Required tools:
- [Golang Migrate](https://github.com/golang-migrate/migrate)
###### How to use

- Add new migration
    ```bash
    migrate create -ext sql -dir db/migrations example_table
    ```

- Run Migration
    - up
      ```bash
      go run cmd/cli/main.go migrate up
      ```
    - down
      ```bash
      go run cmd/cli/main.go migrate down
      ```
more info read the [docs](https://pkg.go.dev/github.com/golang-migrate/migrate/v4).

#### Run
Local run: 

`./script/run-web.sh`

1. Regenerate API Spec
2. Static Check
3. Run Tests
4. Run 

#### Downloads
- [Bundle v0.0.1-dev](https://github.com/aasumitro/pokewar/releases/download/0.0.1-dev/pokewar-bundle.zip)

