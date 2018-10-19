package main

import (
	"flag"
	"io"
	"log"
	"nail/backend/controller"
	"nail/backend/helper"
	"nail/backend/setting"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/tylerb/graceful"
	goji "goji.io"
	"goji.io/pat"
)

// init is for setting max number of procedures based on number of active CPUs
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// main is for setting up environments and start project
func main() {
	// sysExit()
	var confFile string
	flag.StringVar(&confFile, "c", "", "config file")
	flag.Parse()

	//initialize env settings and read from env
	setting.EnvInit()

	if err := mongo.Startup(); err != nil {
		log.Fatal("Mongo startup failed")
		os.Exit(1)
	}

	mux := goji.NewMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   ([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"}),
		AllowedHeaders:   ([]string{"Origin", "Authorization", "Content-Type"}),
		ExposedHeaders:   ([]string{""}),
		MaxAge:           10,
		AllowCredentials: true,
		// OptionsPassthrough: true,
	})
	mux.Use(c.Handler)
	mux.HandleFunc(pat.Get("/"), Root)

	// tenants request handler
	tenantRequest := goji.NewMux()
	tenantRequest.Use(helper.JSON)
	tenantController := controller.TenantController{Name: "TenantController"}
	tenantRequest.HandleFunc(pat.Post("/tenant/"), tenantController.Create)
	tenantRequest.HandleFunc(pat.Get("/tenant/:id"), tenantController.Find)

	// tenants request handler
	accountRequest := goji.NewMux()
	accountRequest.Use(helper.JSON)
	accountController := controller.AccountController{Name: "AccountController"}
	accountRequest.HandleFunc(pat.Get("/account/:id"), accountController.Find)
	accountRequest.HandleFunc(pat.Get("/account/"), accountController.FindAll)
	accountRequest.HandleFunc(pat.Post("/account/"), accountController.Create)
	accountRequest.HandleFunc(pat.Post("/account/:id"), accountController.Update)
	accountRequest.HandleFunc(pat.Delete("/account/"), accountController.Delete)

	// staffs request handler
	staffRequest := goji.NewMux()
	staffRequest.Use(helper.JSON)
	staffController := controller.StaffController{Name: "StaffController"}
	staffRequest.HandleFunc(pat.Get("/staff/:id"), staffController.Find)
	staffRequest.HandleFunc(pat.Get("/staff/"), staffController.FindAll)
	staffRequest.HandleFunc(pat.Post("/staff/"), staffController.Create)
	staffRequest.HandleFunc(pat.Post("/staff/:id"), staffController.Update)
	staffRequest.HandleFunc(pat.Delete("/staff/:id"), staffController.Delete)

	// customers request handler
	customerRequest := goji.NewMux()
	customerRequest.Use(helper.JSON)
	customerController := controller.CustomerController{Name: "CustomerController"}
	customerRequest.HandleFunc(pat.Get("/customer/:id"), customerController.Find)
	customerRequest.HandleFunc(pat.Get("/customer/"), customerController.FindAll)
	customerRequest.HandleFunc(pat.Post("/customer/"), customerController.Create)
	customerRequest.HandleFunc(pat.Post("/customer/:id"), customerController.Update)
	customerRequest.HandleFunc(pat.Delete("/customer/:id"), customerController.Delete)

	// services request handler
	serviceRequest := goji.NewMux()
	serviceRequest.Use(helper.JSON)
	serviceController := controller.ServiceController{Name: "ServiceController"}
	serviceRequest.HandleFunc(pat.Get("/service/:id"), serviceController.Find)
	serviceRequest.HandleFunc(pat.Get("/service/"), serviceController.FindAll)
	serviceRequest.HandleFunc(pat.Post("/service/"), serviceController.Create)
	serviceRequest.HandleFunc(pat.Post("/service/:id"), serviceController.Update)
	serviceRequest.HandleFunc(pat.Delete("/service/:id"), serviceController.Delete)

	// products request handler
	productRequest := goji.NewMux()
	productRequest.Use(helper.JSON)
	productController := controller.ProductController{Name: "ProductController"}
	productRequest.HandleFunc(pat.Get("/product/:id"), productController.Find)
	productRequest.HandleFunc(pat.Get("/product/find_by_service/:id"), productController.FindByService)
	productRequest.HandleFunc(pat.Get("/product/"), productController.FindAll)
	productRequest.HandleFunc(pat.Post("/product/"), productController.Create)
	productRequest.HandleFunc(pat.Post("/product/:id"), productController.Update)
	productRequest.HandleFunc(pat.Delete("/product/:id"), productController.Delete)

	// orders request handler
	orderRequest := goji.NewMux()
	orderRequest.Use(helper.JSON)
	orderController := controller.OrderController{Name: "OrderController"}
	orderRequest.HandleFunc(pat.Get("/order/:id"), orderController.Find)
	orderRequest.HandleFunc(pat.Get("/order/"), orderController.FindAll)
	orderRequest.HandleFunc(pat.Post("/order/"), orderController.Create)
	orderRequest.HandleFunc(pat.Post("/order/:id"), orderController.Update)
	orderRequest.HandleFunc(pat.Post("/order/:id/checkout"), orderController.Checkout)
	orderRequest.HandleFunc(pat.Delete("/order/:id"), orderController.Delete)

	// billings request handler
	billingRequest := goji.NewMux()
	billingRequest.Use(helper.JSON)
	billingController := controller.BillingController{Name: "BillingController"}
	billingRequest.HandleFunc(pat.Get("/billing/:id"), billingController.Find)
	billingRequest.HandleFunc(pat.Get("/billing/"), billingController.FindAll)
	billingRequest.HandleFunc(pat.Post("/billing/"), billingController.Create)
	billingRequest.HandleFunc(pat.Post("/billing/:id"), billingController.Update)
	billingRequest.HandleFunc(pat.Delete("/billing/:id"), billingController.Delete)

	// sessions request handler
	sessionRequest := goji.NewMux()
	sessionRequest.Use(helper.JSON)
	sessionController := controller.SessionController{Name: "SessionController"}
	sessionRequest.HandleFunc(pat.Post("/session/login"), sessionController.Login)
	sessionRequest.HandleFunc(pat.Post("/session/logout"), sessionController.Logout)

	mux.Handle(pat.New("/tenant/*"), tenantRequest)
	mux.Handle(pat.New("/account/*"), accountRequest)
	mux.Handle(pat.New("/staff/*"), staffRequest)
	mux.Handle(pat.New("/customer/*"), customerRequest)
	mux.Handle(pat.New("/service/*"), serviceRequest)
	mux.Handle(pat.New("/product/*"), productRequest)
	mux.Handle(pat.New("/order/*"), orderRequest)
	mux.Handle(pat.New("/billing/*"), billingRequest)
	mux.Handle(pat.New("/session/*"), sessionRequest)

	srv := &graceful.Server{
		Timeout: 3 * time.Second,
		BeforeShutdown: func() bool {
			cleanup()
			return true
		},
		Server: &http.Server{
			Addr:    ":4011",
			Handler: mux,
		},
	}
	log.Println("Server starts and listens at port 4011")
	srv.ListenAndServe()
}

// Root displays orginial name of service
func Root(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Backend\nPowered by LAI\n=============================\n")
}

func sysExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

func cleanup() {
	log.Println("Application termination request received")
	log.Println("shutting down mongo")
	mongo.Shutdown()
	log.Println("Terminating application")
}
