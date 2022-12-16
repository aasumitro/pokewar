package sql

type rankSQLRepository struct {
	//db *sql.DB
}

func NewRankSQLRepository() any {
	return rankSQLRepository{}
}
