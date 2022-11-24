package types

import (
	"moonspace/model"
	repo_types "moonspace/repository/types"

	"github.com/qiniu/qmgo"
)

type MongoRepository[T model.Entity] interface {
	repo_types.Repository[T]
	SwitchTable(data *T) (*qmgo.Collection, error)
}
