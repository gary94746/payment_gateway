package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"payment-processor.gary94746/main/lib/processors"
)

func GetPayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.services.GetPayment(paymentId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, payment)
}

func CapturePayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")

	errors := api.services.CapturePayment(paymentId)
	if errors != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errors.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func CreatePayment(ctx *gin.Context) {
	var body processors.Payment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := api.services.CreatePayment(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	}

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
	refund, err := api.services.RefundPayment(paymentId, body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": refund,
	})
}
