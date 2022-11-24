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

type ProdRequest struct {
	ProductID  uint64 `json:"product_id"`
	CategoryID uint64 `json:"category_id"`
}

func CreateProduct(s service.Product, serverUrl string, ce utils.ClaimsExtractor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createProduct(ctx, s, serverUrl, ce)
	}
}

func DeleteProduct(s service.Product) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getDeleteProduct(ctx, s)
	}
}

func UpdateProduct(s service.Product) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		updateProduct(ctx, s)
	}
}

func GetProduct(s service.Product) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getDeleteProduct(ctx, s)
	}
}

func GetProductsLimit(s service.Product) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		getLimitProducts(ctx, s)
	}
}

func createProduct(ctx *gin.Context, s service.Product, serverUrl string, ce utils.ClaimsExtractor) {
	token := ctx.Request.Header.Get(utils.TokenHeader)
	uid, err := ce.Extract(token, utils.OAuthClaimSubject)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	p := dto.ProductDto{}
	err = ctx.Bind(p)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	imgPath, err := utils.SaveImage(ctx, p.Image)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cid := ctx.Param("categoryId")
	prod := dto.ProductDtoToModel(p, serverUrl+"/"+imgPath, cid, uid)
	err = s.Create(prod)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, "OK")
}

func updateProduct(ctx *gin.Context, s service.Product) {
	p := model.Product{}
	err := utils.DecodeRequestBody(ctx.Request.Body, &p)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	cid := ctx.Param("categoryId")
	pid := ctx.Param("productId")
	p.CategoryID = cid
	p.ID = pid
	err = s.Update(p)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(200, "OK")
}

func getDeleteProduct(ctx *gin.Context, service service.Product) {
	cid := ctx.Param("categoryId")
	pid := ctx.Param("productId")
	var res any
	var err error

	if ctx.Request.Method == http.MethodGet {
		res, err = service.Get(cid, pid)
	} else if ctx.Request.Method == http.MethodDelete {
		err = service.Delete(cid, pid)
		res = "OK"
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, res)
}

func getLimitProducts(ctx *gin.Context, s service.Product) {
	cid := ctx.Param("categoryId")
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

	res, err := s.GetProductsLimit(cid, start, end)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, res)
}
