package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"payment-processor.gary94746/main/app/services"
	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

type Api struct {
	database database.Database
	services services.Services
}

var api Api

func main() {
	godotenv.Load()

	paypal := &processors.PayPal{}
	stripe := &processors.Stripe{}
	inMemory := database.InMemory{}

	stripe.Init(processors.PaymentSettings{
		Credentials: map[string]string{
			"token": os.Getenv("STRIPE_TOKEN"),
		},
	})
	paypal.Init(processors.PaymentSettings{
		Credentials: map[string]string{
			"client_id":    os.Getenv("PAYPAL_CLIENT_ID"),
			"client_token": os.Getenv("PAYPAL_CLIENT_TOKEN"),
			"mode":         os.Getenv("PAYPAL_MODE"),
		},
	})

	api = Api{
		database: inMemory,
		services: services.Services{
			Database:         inMemory,
			PaymentProcessor: paypal,
		},
	}

	r := gin.Default()

	processorV1Group := r.Group("v1/processor/payment")
	processorV1Group.GET("/:id", GetPayment)
	processorV1Group.POST("/", CreatePayment)
	processorV1Group.POST("/:id/capture", CapturePayment)
	processorV1Group.POST("/:id/refund", RefundPayment)

	r.Run("127.0.0.1:3000")
}
