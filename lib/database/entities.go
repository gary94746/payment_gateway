package database

type PartialRefund struct {
	Amount int64 `json:"amount"`
}

type RefundResponse struct {
	Id     string `json:"id"`
	Amount string `json:"amount"`
}

type LineItem struct {
	Name     string `json:"name"`
	Amount   int64  `json:"amount"`
	Quantity int32  `json:"quantity"`
}

type Payment struct {
	Currency    string           `json:"currency"`
	Amount      int64            `json:"amount"`
	Status      string           `json:"status"`
	RedirectUrl string           `json:"redirectUrl"`
	CancelUrl   string           `json:"cancelUrl"`
	PrivateId   string           `json:"privateId"`
	LineItems   []LineItem       `json:"lineItems"`
	Refunds     []RefundResponse `json:"refunds"`
	Id          string           `json:"id"`
}

type Database interface {
	Save(payment Payment) string
	FindById(id string) (*Payment, error)
	UpdateStatus(id string, status string) error
	AttachRefund(paymentId string, refund RefundResponse) error
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
