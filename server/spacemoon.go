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
	log.Default().Print("v0.2.1")
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("registering server handlers...")

	http.Handle("/login", login.NewHandler(loginPersistence, time.Hour))
	protector := login.NewProtector(loginPersistence)

	productHandler := product_handler.MakeHandler(&temporaryProductPersistence{})
	preparedProductHandler := prepareHandler(protector, productHandler, http.MethodGet)
	http.Handle("/product", preparedProductHandler)

	categoryHandler := category_handler.MakeHandler(&temporaryCategoryPersistence{})
	preparedCategoryHandler := prepareHandler(protector, categoryHandler, http.MethodGet)
	http.Handle("/category", preparedCategoryHandler)

	log.Default().Print("handler registration done, ready for takeoff")
}

func prepareHandler(protector login.Protector, handler http.Handler, methods ...string) http.Handler {
	protectedHandler := protector.Protect(&handler)
	for _, method := range methods {
		protectedHandler.Unprotect(method)
	}
	corsEnabledProtectedProductHandler := cors.EnableCors(protectedHandler)
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
