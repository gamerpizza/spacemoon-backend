package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"spacemoon/login"
	"spacemoon/server/category_handler"
	"spacemoon/server/cors"
	"spacemoon/server/network_handler"
	"spacemoon/server/product_handler"
	"strings"
	"time"
)

func main() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("v0.6.2")
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("registering server handlers...")

	loginPersistence := getLoginPersistence()
	corsEnabledLoginHandler := cors.EnableCors(login.NewHandler(loginPersistence, time.Hour), http.MethodGet, http.MethodPost)
	http.Handle("/login", corsEnabledLoginHandler)

	protector := login.NewProtector(loginPersistence)

	socialNetworkHandler := network_handler.New(getSocialNetworkPersistence(), loginPersistence)
	corsEnabledSocialNetworkHandler := cors.EnableCors(socialNetworkHandler, http.MethodGet, http.MethodPost)
	protectedSocialNetworkHandler := protector.Protect(&corsEnabledSocialNetworkHandler)
	protectedSocialNetworkHandler.Unprotect(http.MethodGet)
	http.Handle("/posts", protectedSocialNetworkHandler)

	productHandler := product_handler.MakeHandler(getProductPersistence(), loginPersistence)
	preparedProductHandler := prepareHandler(protector, productHandler, http.MethodGet)
	http.Handle("/product", preparedProductHandler)
	productRatingHandler := product_handler.MakeRankingsHandler(getProductRatingsPersistence())
	preparedProductRatingHandler := prepareHandler(protector, productRatingHandler, http.MethodGet)
	http.Handle("/product/rating", preparedProductRatingHandler)

	categoryHandler := category_handler.MakeHandler(getCategoryPersistence())
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
	certFile := os.Getenv("cert_file")
	keyFile := os.Getenv("key_file")
	if strings.TrimSpace(certFile) == "" || strings.TrimSpace(keyFile) == "" {
		log.Default().Print("serving without TLS\n")
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
		}
		return
	}
	log.Default().Print("serving using TLS\n")
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, nil)
	if err != nil {
		log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
	}
}

const port = 1234
