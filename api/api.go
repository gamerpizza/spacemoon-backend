package api

import (
	"fmt"
	"moonspace/api/handlers"
	"moonspace/repository"
	"moonspace/service"
	"moonspace/utils"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v11"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	apiUrl = "/api"

	product          = "/product"
	createProduct    = "/:categoryId/create"
	deleteProduct    = "/:categoryId/:productId/delete"
	updateProduct    = "/:categoryId/:productId/update"
	getProduct       = "/:categoryId/:productId/get"
	paginateProducts = "/:categoryId/get/:start/:end"

	category           = "/category"
	createCategory     = "/create"
	deleteCategory     = "/:id/delete"
	updateCategory     = "/:id/update"
	getCategory        = "/:id/get"
	paginateCategories = "/get/:start/:end"

	order                = "/order"
	createOrderPaypal    = "/paypal/createOrder"
	createOrderStripe    = "/stripe/createOrder"
	createOrderGooglePay = "/googlePay/createOrder"
	checkoutPaypal       = "/paypal/:orderId/checkout"
	checkoutStripe       = "/stripe/:orderId/checkout"

	images   = "/images"
	getImage = "/:imageName"
)

type API struct {
	cfg        *Config
	service    *service.Service
	router     *gin.Engine
	authServer *gocloak.GoCloak
}

func New(cfg *Config) *API {
	gin.DefaultWriter = os.Stdout
	gin.DefaultErrorWriter = os.Stdout
	cli := repository.MakeRepositoryClient(*cfg.DB)

	return &API{
		cfg:        cfg,
		service:    service.NewService(cli, *cfg.DB, *cfg.Payment),
		router:     gin.Default(),
		authServer: nil,
	}
}

func (a *API) Start() {
	if a.cfg.Security != nil {
		if a.cfg.Security.CORS != nil {
			setupCors(a.router, *a.cfg.Security.CORS)
		}
		if a.cfg.Security.OAuth != nil {
			keycloak := gocloak.NewClient(a.cfg.Security.OAuth.URL, gocloak.SetAuthRealms("realms"), gocloak.SetAuthAdminRealms("admin/realms"))
			a.authServer = &keycloak
		}

		a.setupSecureRoutes(a.router, *a.service)
	}

	a.router.Run(fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port))
}

func (a *API) setupSecureRoutes(router *gin.Engine, service service.Service) *gin.RouterGroup {
	extractor := utils.NewClaimsExtractorKeycloak(*a.authServer, a.cfg.Security.OAuth.Realm)

	apiRouter := router.Group(apiUrl)

	// These should be enabled per router group.
	// Router groups should be grouped by roles.

	// apiRouter.Use(middleware.Authentication(*a.authServer, a.cfg.Security.OAuth.Realm))
	// testRole := []string{middleware.RoleAdmin}
	// apiRouter.Use(middleware.AuthorizationForRoles(a.cfg.Security.OAuth.AdminData, testRole)(*a.authServer, a.cfg.Security.OAuth.Realm, a.cfg.Security.OAuth.ClientID))

	url := "http://" + a.cfg.Host + ":" + fmt.Sprintf("%d", a.cfg.Port)
	categoryRouter := apiRouter.Group(category)
	categoryRouter.POST(createCategory, handlers.CreateCategory(service.Category, url, &extractor))
	categoryRouter.DELETE(deleteCategory, handlers.DeleteCategory(service.Category))
	categoryRouter.PUT(updateCategory, handlers.UpdateCategory(service.Category))
	categoryRouter.GET(getCategory, handlers.GetCategory(service.Category))
	categoryRouter.GET(paginateCategories, handlers.GetCategoryLimit(service.Category))

	productRouter := apiRouter.Group(product)
	productRouter.POST(createProduct, handlers.CreateProduct(service.Product, url, &extractor))
	productRouter.DELETE(deleteProduct, handlers.DeleteProduct(service.Product))
	productRouter.PUT(updateProduct, handlers.UpdateProduct(service.Product))
	productRouter.GET(getProduct, handlers.GetProduct(service.Product))
	productRouter.GET(paginateProducts, handlers.GetProductsLimit(service.Product))

	orderRouter := apiRouter.Group(order)
	orderRouter.POST(createOrderPaypal, handlers.CreatePayPalOrder(service.Order))
	orderRouter.POST(createOrderStripe, handlers.CreateStripeOrder(service.Order))
	orderRouter.POST(createOrderGooglePay, handlers.CreateGooglePayOrder(service.Order))
	orderRouter.POST(checkoutPaypal, handlers.PayPalCheckout(service.Order))
	orderRouter.POST(checkoutStripe, handlers.StripeCheckout(service.Order))

	imagesRouter := router.Group(images)
	imagesRouter.GET(getImage, handlers.GetImage(url))

	return apiRouter
}

func setupCors(r gin.IRoutes, cfg CORS) gin.IRoutes {
	return r.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.AllowOrigins,
		AllowMethods:  cfg.AllowMethods,
		AllowHeaders:  cfg.AllowHeaders,
		ExposeHeaders: cfg.ExposeHeaders,
		MaxAge:        time.Hour,
	}))
}
