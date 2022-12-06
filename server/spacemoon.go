package main

import (
	"fmt"
	"log"
	"net/http"
	"spacemoon/login"
	"spacemoon/server/category_handler"
	"spacemoon/server/product_handler"
	"time"
)

func main() {
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("registering server handlers...")
	http.Handle("/product", product_handler.MakeHandler(&temporaryProductPersistence{}))
	http.Handle("/category", category_handler.MakeHandler(&temporaryCategoryPersistence{}))
	http.Handle("/login", login.NewHandler(loginPersistence, time.Hour))
	log.Default().Print("handler registration done, ready for takeoff")
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
