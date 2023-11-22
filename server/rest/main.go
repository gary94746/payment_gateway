package rest

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"payment-processor.gary94746/main/app/services"
	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

type ApiRest struct {
	database database.Database
	services services.Services
}

func (ar ApiRest) Serve() error {
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

	api := ApiRest{
		database: inMemory,
		services: services.Services{
			Database:         inMemory,
			PaymentProcessor: paypal,
		},
	}

	r := gin.Default()

	r.GET("/api/health", health)

	processorV1Group := r.Group("/api/v1/processor/payment")
	processorV1Group.GET("/:id", api.getPayment)
	processorV1Group.POST("/", api.createPayment)
	processorV1Group.POST("/:id/capture", api.capturePayment)
	processorV1Group.POST("/:id/refund", api.refundPayment)

	error := r.Run("127.0.0.1:3000")
	return error
}

func health(ctx *gin.Context) {
	response := struct{ success bool }{
		success: true,
	}

	ctx.JSON(http.StatusOK, response)
}
