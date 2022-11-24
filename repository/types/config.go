package types

type DatabaseType string

const (
	Mongo    DatabaseType = "mongo"
	Postgres DatabaseType = "postgres"
)

type RepositoryType string

const (
	Cart    RepositoryType = "cart"
	Product RepositoryType = "product"
	Order   RepositoryType = "order"
)

type Config struct {
	Url      string       `yaml:"url"`
	Database string       `yaml:"database"`
	Type     DatabaseType `yaml:"type"`
}
