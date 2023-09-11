package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

func GetPayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.Storage.FindById(paymentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, payment)
}

func CapturePayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.Storage.FindById(paymentId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "payment not found",
		})
		return
	}

	_, captureErr := api.PaymentProcessor.Capture(payment.PrivateId)
	if captureErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": captureErr.Error(),
		})

		return
	}

	api.Storage.UpdateStatus(paymentId, processors.StatusCaptured)
	ctx.JSON(http.StatusOK, gin.H{})
}

func CreatePayment(ctx *gin.Context) {
	var body processors.Payment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := api.PaymentProcessor.Create(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	items := []database.LineItem{}

	for _, item := range body.LineItems {
		items = append(items, database.LineItem{
			Name:     item.Name,
			Amount:   item.Amount,
			Quantity: item.Quantity,
		})
	}

	databasePayment := database.Payment{
		Currency:    body.Currency,
		Amount:      body.Amount,
		Status:      processors.StatusCreated,
		RedirectUrl: body.RedirectUrl,
		CancelUrl:   body.CancelUrl,
		PrivateId:   payment.PrivateId,
		Id:          payment.Id,
		Reference:   body.Reference,
		LineItems:   items,
		Refunds:     []database.RefundResponse{},
		Customer:    database.Customer(body.Customer),
	}

	paymentId := api.Storage.Save(databasePayment)
	payment.Id = paymentId

	ctx.JSON(http.StatusOK, gin.H{
		"data": payment,
	})
}

func RefundPayment(ctx *gin.Context) {
	var body processors.PartialRefund
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentId, _ := ctx.Params.Get("id")
	order, err := api.Storage.FindById(paymentId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	refund, err := api.PaymentProcessor.Refund(order.PrivateId, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	api.Storage.UpdateStatus(paymentId, processors.StatusRefunded)
	api.Storage.AttachRefund(paymentId, database.RefundResponse{
		Id:     refund.Id,
		Amount: refund.Amount,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"data": refund,
	})
}
