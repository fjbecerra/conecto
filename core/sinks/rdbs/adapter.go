package rdbs

type Adapter interface {
	DriverName() string
	Placeholder(n int) string
	BuildInsert(table string, columns []string, rows int) string
}