package main

import (
	"net/http"
	"spacemoon/server/product"
)

func main() {
	http.Handle("/product", product.Handler{})
}
