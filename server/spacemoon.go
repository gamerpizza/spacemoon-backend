package main

import (
	"fmt"
	"log"
	"net/http"
	"spacemoon/login"
	"spacemoon/server/category_handler"
	"spacemoon/server/cors"
	"spacemoon/server/product_handler"
	"time"
)

func main() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("v0.2.3")
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("registering server handlers...")

	corsEnabledLoginHandler := cors.EnableCors(login.NewHandler(loginPersistence, time.Hour), http.MethodGet)
	http.Handle("/login", corsEnabledLoginHandler)
	protector := login.NewProtector(loginPersistence)

	productHandler := product_handler.MakeHandler(&temporaryProductPersistence{})
	preparedProductHandler := prepareHandler(protector, productHandler, http.MethodGet)
	http.Handle("/product", preparedProductHandler)
	productRatingHandler := product_handler.MakeRankingsHandler()
	preparedProductRatingHandler := prepareHandler(protector, productRatingHandler, http.MethodGet)
	http.Handle("/product/rating", preparedProductRatingHandler)

	categoryHandler := category_handler.MakeHandler(&temporaryCategoryPersistence{})
	preparedCategoryHandler := prepareHandler(protector, categoryHandler, http.MethodGet)
	http.Handle("/category", preparedCategoryHandler)

	log.Default().Print("handler registration done, ready for takeoff")
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

func listenAndServe() {
	log.Default().Print(getRandomSpaceQuote())
	log.Default().Printf("listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
	}
}

const port = 1234
