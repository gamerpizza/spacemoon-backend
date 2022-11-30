package handlers

import (
	"fmt"
	"moonspace/dto"
	"moonspace/model"
	"moonspace/service"
	"moonspace/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryAction func(cat model.Category) error

func CreateCategory(s service.Category, serverUrl string, ce utils.ClaimsExtractor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createCategory(ctx, s, serverUrl, ce)
	}
}

func GetCategory(s service.Category) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getCategory(ctx, s)
	}
}

func GetCategoryLimit(s service.Category) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getCategoryLimit(ctx, s)
	}
}

func UpdateCategory(s service.Category, serverUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		updateCategory(ctx, serverUrl, s)
	}
}

func DeleteCategory(s service.Category) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deleteCategory(ctx, s)
	}
}

func createCategory(ctx *gin.Context, s service.Category, serverUrl string, ce utils.ClaimsExtractor) {
	token := ctx.Request.Header.Get(utils.TokenHeader)
	uid, err := ce.Extract(token, utils.OAuthClaimSubject)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c := dto.CategoryDto{}
	err = ctx.Bind(&c)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	imgPath, err := utils.SaveImage(ctx, c.Image)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	category := dto.CategoryDtoToModel(c, serverUrl+"/"+imgPath, uid)
	err = s.Create(category)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, "OK")
}

func getCategory(ctx *gin.Context, s service.Category) {
	id := ctx.Param("id")
	cat, err := s.Get(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, cat)
}

func getCategoryLimit(ctx *gin.Context, s service.Category) {
	start, err := utils.StringToUint64(ctx.Param("start"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	end, err := utils.StringToUint64(ctx.Param("end"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if start > end {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("start > end"))
		return
	}

	cats, err := s.GetLimit(start, end)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, cats)
}

func updateCategory(ctx *gin.Context, serverUrl string, s service.Category) {
	c := dto.CategoryDto{}
	err := ctx.Bind(&c)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	imgPath, err := utils.SaveImage(ctx, c.Image)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cat := model.Category{}
	id := ctx.Param("id")
	cat.Name = c.Name
	cat.Image = serverUrl + "/" + imgPath
	err = s.Update(id, cat)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, "OK")
}

func deleteCategory(ctx *gin.Context, s service.Category) {
	id := ctx.Param("id")
	err := s.Delete(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(200, "OK")
}
