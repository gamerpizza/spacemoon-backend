package handlers_test

import (
	"encoding/json"
	"fmt"
	"moonspace/api/handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetSchema(t *testing.T) {
	r := gin.Default()
	r.GET("/schema/:name", handlers.GetMetaschema())

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/schema/category", nil)
	assert.NoError(t, err)

	r.ServeHTTP(w, req)

	var result interface{}
	assert.Equal(t, 200, w.Result().StatusCode)
	err = json.Unmarshal(w.Body.Bytes(), &result)
	fmt.Println("RESULT: ", result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
