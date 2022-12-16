package sql

type battleSQLRepository struct {
	//db *sql.DB
}

func NewBattleSQLRepository() any {
	return &battleSQLRepository{}
}
