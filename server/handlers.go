package main

import (
	"net/http"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/product/category"
	"spacemoon/product/handler"
	"spacemoon/product/ratings"
	"spacemoon/server/cors"
)

func setupSocialNetworkHandler(snp network.Persistence, loginPersistence login.Persistence, mediaFilePersistence network.MediaFilePersistence, protector login.Protector) {
	socialNetworkHandler := network.NewHandler(snp, loginPersistence, mediaFilePersistence)
	protectedSocialNetworkHandler := protector.Protect(&socialNetworkHandler)
	protectedSocialNetworkHandler.Unprotect(http.MethodGet)
	corsEnabledSocialNetworkHandler := cors.EnableCors(protectedSocialNetworkHandler,
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
	http.Handle("/posts", corsEnabledSocialNetworkHandler)
}

func setupProductHandlers(loginPersistence login.Persistence, protector login.Protector) {
	productHandler := handler.MakeHandler(getProductPersistence(), loginPersistence)
	preparedProductHandler := prepareHandler(protector, productHandler, http.MethodGet)
	http.Handle("/product", preparedProductHandler)
	productRatingHandler := ratings.MakeRankingsHandler(getProductRatingsPersistence())
	preparedProductRatingHandler := prepareHandler(protector, productRatingHandler, http.MethodGet)
	http.Handle("/product/rating", preparedProductRatingHandler)
}

func setupCategoryHandler(protector login.Protector) {
	categoryHandler := category.MakeHandler(getCategoryPersistence())
	preparedCategoryHandler := prepareHandler(protector, categoryHandler, http.MethodGet)
	http.Handle("/category", preparedCategoryHandler)
}

func prepareHandler(protector login.Protector, handler http.Handler, unprotectedMethods ...string) http.Handler {
	protectedHandler := protector.Protect(&handler)
	for _, method := range unprotectedMethods {
		protectedHandler.Unprotect(method)
	}
	corsEnabledProtectedProductHandler := cors.EnableCors(protectedHandler,
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
	return corsEnabledProtectedProductHandler
}
