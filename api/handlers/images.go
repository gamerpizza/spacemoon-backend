package handlers

import (
	"fmt"
	"moonspace/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetImage(serverUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		imgName := strings.ReplaceAll(ctx.Param("imageName"), serverUrl, "")
		imgPath := utils.ImagesPath
		fullImgPath := imgPath + imgName
		_, err := os.Stat(fullImgPath)

		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}

		if os.IsNotExist(err) {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("Image %s does not exist.", fullImgPath))
			return
		}

		ctx.File(fullImgPath)
	}
}
