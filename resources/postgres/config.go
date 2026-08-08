package postgres

type PostgresResourceConfig struct{}

type PostgresSinkConfig struct {
	Table string `json:"table"`
	AutoCreate bool `json:"auto_create"`
}
