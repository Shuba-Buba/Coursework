package storage

type Dao interface {
	Append(table string, data string)
	GetRows(table string) chan string
	GetAllTables() []string
}
