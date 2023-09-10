package main

type CheckoutResponse struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}
type SessionResponse struct {
	Id            string `json:"id"`
	PaymentStatus string `json:"payment_status"`
}
