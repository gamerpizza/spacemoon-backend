package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/orderedmap"
)

const schemasFolder = "../../schemas/"

func GetMetaschema() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getMetaschema(ctx)
	}
}

func getMetaschema(ctx *gin.Context) {
	name := ctx.Param("name")
	filename := schemasFolder + name + ".json"

	if _, err := os.Stat(filename); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("File does not exist"))
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error reading file")
		return
	}

	schema := orderedmap.OrderedMap{}
	err = schema.UnmarshalJSON(bytes)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Error parsing json"))
	}

	ctx.JSON(http.StatusOK, schema)
}
