package processors

type CheckoutResponse struct {
	Url           string `json:"url"`
	Id            string `json:"id"`
	PaymentIntent string `json:"payment_intent"`
}
type PaymentIntentResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
