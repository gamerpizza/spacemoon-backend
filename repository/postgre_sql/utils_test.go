package postgre_sql_test

import (
	"moonspace/model"
	"moonspace/repository/postgre_sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Query(t *testing.T) {
	c := model.Cart{
		UserID: "ASDF",
	}

	query, vals, err := postgre_sql.CreateQueryFromkey(c.Key())
	expectedQuery := "id = ?"
	expectedVals := []any{uint64(1)}

	assert.NoError(t, err)
	assert.Equal(t, expectedQuery, query)
	assert.EqualValues(t, expectedVals, vals)
}
