package main

import (
	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

type Api struct {
	Storage          database.Storage
	PaymentProcessor processors.PaymentConnector
}
