package types

import (
	"moonspace/model"
	repo_types "moonspace/repository/types"
)

type PostgreRepository[T model.Entity] interface {
	repo_types.Repository[T]
}
