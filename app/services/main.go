package services

import (
	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

type Services struct {
	Database         database.Database
	PaymentProcessor processors.PaymentConnector
}
