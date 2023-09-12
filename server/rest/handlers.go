package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"payment-processor.gary94746/main/lib/processors"
)

func (api ApiRest) getPayment(ctx *gin.Context) {
	paymentId := ctx.Param("id")
	payment, err := api.services.GetPayment(paymentId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, payment)
}

func (api ApiRest) capturePayment(ctx *gin.Context) {
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

func (api ApiRest) createPayment(ctx *gin.Context) {
	var body Payment
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var items []processors.LineItem
	for _, item := range body.LineItems {
		items = append(items, processors.LineItem{
			Name:     item.Name,
			Amount:   item.Amount,
			Quantity: item.Quantity,
		})
	}

	paymentPayload := processors.Payment{
		Currency:    body.Currency,
		Amount:      body.Amount,
		RedirectUrl: body.RedirectUrl,
		CancelUrl:   body.CancelUrl,
		LineItems:   items,
	}

	payment, err := api.services.CreatePayment(paymentPayload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": payment,
	})
}

func (api ApiRest) refundPayment(ctx *gin.Context) {
	var body PartialRefund
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refundPayload := processors.PartialRefund{
		Amount: body.Amount,
	}
	paymentId, _ := ctx.Params.Get("id")
	refund, err := api.services.RefundPayment(paymentId, refundPayload)

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
