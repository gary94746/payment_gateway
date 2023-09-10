package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var api Api

func main() {
	godotenv.Load()

	paypal := &PayPal{}
	stripe := &Stripe{}
	inMemory := InMemory{}

	stripe.Init(PaymentSettings{
		Credentials: map[string]string{
			"token": os.Getenv("STRIPE_TOKEN"),
		},
	})
	paypal.Init(PaymentSettings{
		Credentials: map[string]string{
			"client_id":    os.Getenv("PAYPAL_CLIENT_ID"),
			"client_token": os.Getenv("PAYPAL_CLIENT_TOKEN"),
			"mode":         os.Getenv("PAYPAL_MODE"),
		},
	})

	api = Api{
		PaymentProcessor: paypal,
		Storage:          inMemory,
	}

	r := gin.Default()

	processorV1Group := r.Group("v1/processor/payment")
	processorV1Group.GET("/:id", GetPayment)
	processorV1Group.POST("/", CreatePayment)
	processorV1Group.POST("/:id/capture", CapturePayment)
	processorV1Group.POST("/:id/refund", RefundPayment)

	r.Run("127.0.0.1:3000")
}
