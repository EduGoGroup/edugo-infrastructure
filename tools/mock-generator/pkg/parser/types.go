package parser

// TableData contiene los datos extraidos de un INSERT
type TableData struct {
	Table   string
	Columns []string
	Rows    [][]interface{}
}
