package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.Storage.findById(paymentId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, payment)
}

func CapturePayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.Storage.findById(paymentId)
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

	api.Storage.updateStatus(paymentId, "CAPTURED")
	ctx.JSON(http.StatusOK, gin.H{})
}

func CreatePayment(ctx *gin.Context) {
	var body Payment
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

	body.PrivateId = payment.PrivateId
	body.Status = "CREATED"
	paymentId := api.Storage.save(body)
	payment.Id = paymentId

	ctx.JSON(http.StatusOK, gin.H{
		"data": payment,
	})
}

func RefundPayment(ctx *gin.Context) {
	var body PartialRefund
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentId, _ := ctx.Params.Get("id")
	order, err := api.Storage.findById(paymentId)
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

	api.Storage.updateStatus(paymentId, "REFUNDED")
	api.Storage.attachRefund(paymentId, *refund)

	ctx.JSON(http.StatusOK, gin.H{
		"data": refund,
	})
}
