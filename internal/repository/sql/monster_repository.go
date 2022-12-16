package sql

type monsterSQLRepository struct {
	//db *sql.DB
}

func NewMonsterSQlRepository() any {
	return &monsterSQLRepository{}
}
