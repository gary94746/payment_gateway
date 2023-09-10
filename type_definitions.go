package main

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

type Customer struct {
	Name     string `json:"name" binding:"required,min=2,max=300"`
	LastName string `json:"lastName" binding:"required,min=2,max=300"`
	Email    string `json:"email" binding:"required,email,min=2,max=300"`
}

type Payment struct {
	Currency    string           `json:"currency" binding:"required,iso4217"`
	Amount      int64            `json:"amount" binding:"required,number,min=1000"`
	Status      string           `json:"status" binding:"-"`
	RedirectUrl string           `json:"redirectUrl" binding:"required,url"`
	CancelUrl   string           `json:"cancelUrl" binding:"required,url"`
	PrivateId   string           `json:"privateId" binding:"-"`
	Reference   string           `json:"reference" binding:"required,min=1,max=10"`
	LineItems   []LineItem       `json:"lineItems" binding:"required,gt=0,dive,lt=200,dive"`
	Refunds     []RefundResponse `json:"refunds"`
	Customer    Customer         `json:"customer" binding:"required,dive"`
	Id          string           `json:"id" binding:"-"`
}

type Storage interface {
	save(payment Payment) string
	findById(id string) (*Payment, error)
	updateStatus(id string, status string) error
	attachRefund(paymentId string, refund RefundResponse) error
}

type PaymentDetail struct {
	Id          string `json:"id"`
	PrivateId   string `json:"privateId"`
	RedirectUrl string `json:"redirectUrl"`
	Status      string `json:"status"`
}

type PaymentSettings struct {
	Credentials map[string]string
	Mode        string
}

type PaymentConnector interface {
	Init(settings PaymentSettings) error
	Create(payment Payment) (*PaymentDetail, error)
	Capture(paymentId string) (bool, error)
	Refund(paymentId string, refund PartialRefund) (*RefundResponse, error)
}

type Api struct {
	Storage          Storage
	PaymentProcessor PaymentConnector
}
