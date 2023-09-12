package rest

type PartialRefund struct {
	Amount int64 `json:"amount" binding:"required,number"`
}

type RefundResponse struct {
	Id     string `json:"id"`
	Amount string `json:"amount"`
}

type LineItem struct {
	Name     string `json:"name" binding:"required,min=1,max=400"`
	Amount   int64  `json:"amount" binding:"required,number,min=1000"`
	Quantity int32  `json:"quantity" binding:"required,number,min=1"`
}

type Payment struct {
	Currency    string           `json:"currency" binding:"required,iso4217"`
	Amount      int64            `json:"amount" binding:"required,number,min=1000"`
	Status      string           `json:"status" binding:"-"`
	RedirectUrl string           `json:"redirectUrl" binding:"required,url"`
	CancelUrl   string           `json:"cancelUrl" binding:"required,url"`
	PrivateId   string           `json:"privateId" binding:"-"`
	LineItems   []LineItem       `json:"lineItems" binding:"required,gt=0,dive,lt=200,dive"`
	Refunds     []RefundResponse `json:"refunds"`
	Id          string           `json:"id" binding:"-"`
}

type PaymentDetail struct {
	Id          string `json:"id"`
	PrivateId   string `json:"privateId"`
	RedirectUrl string `json:"redirectUrl"`
	Status      string `json:"status"`
}
