package processors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type PayPal struct {
	Client      *http.Client
	basePath    string
	log         slog.Logger
	username    string
	password    string
	bearerToken string
}

func (p *PayPal) Init(settings PaymentSettings) error {
	p.log = *slog.New(slog.NewJSONHandler(os.Stdout, nil))
	p.username = settings.Credentials["client_id"]
	p.basePath = "https://api.paypal.com"
	p.password = settings.Credentials["client_token"]

	isSandbox := settings.Credentials["mode"] == "SANDBOX"
	if isSandbox {
		p.basePath = "https://api.sandbox.paypal.com"
	}

	p.Client = &http.Client{
		Timeout: 60 * time.Second,
	}

	return nil
}

func (p *PayPal) Create(payment Payment) (*PaymentDetail, error) {
	items := []Item{}

	for _, lineItem := range payment.LineItems {
		item := Item{
			Name:     lineItem.Name,
			Quantity: strconv.Itoa(int(lineItem.Quantity)),
			UnitAmount: Amount{
				CurrencyCode: payment.Currency,
				Value:        strconv.Itoa(int(lineItem.Amount / 100)),
			},
		}

		items = append(items, item)
	}

	purchaseUnit := PurchaseUnits{
		Amount: PurchaseUnitAmount{
			CurrencyCode: payment.Currency,
			Value:        strconv.Itoa(int(payment.Amount / 100)),
			Breakdown: Breakdown{
				ItemTotal: Amount{
					CurrencyCode: payment.Currency,
					Value:        strconv.Itoa(int(payment.Amount / 100)),
				},
			},
		},
		Items: items,
	}

	order := Order{
		Intent: "CAPTURE",
		ApplicationContext: ApplicationContext{
			ReturnUrl: payment.RedirectUrl,
			CancelUrl: payment.CancelUrl,
		},
		PurchaseUnits: []PurchaseUnits{purchaseUnit},
	}

	payload, err := json.Marshal(order)
	if err != nil {
		p.log.Info("Error on marshal order", err)
	}

	request, err := http.NewRequest(http.MethodPost, p.basePath+"/v2/checkout/orders", bytes.NewBuffer(payload))
	if err != nil {
		p.log.Info("Error on request", err)
	}

	response, err := p.requestWrapper(*request)
	if err != nil {
		p.log.Info("Do request err", err)
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		p.log.Error("Error decoding order response", rawResponse)
		return nil, errors.New("error decoding order response")
	}

	defer response.Body.Close()

	isCreatedStatus := response.StatusCode == http.StatusCreated
	if !isCreatedStatus {
		p.log.Error("PAYPAL_ORDER_CREATION_ERROR", "RESPONSE", string(rawResponse))

		return nil, fmt.Errorf("error creating the order, detail -> %s", string(rawResponse))
	}

	orderResponse := &OrderResponse{}
	dErr := json.Unmarshal([]byte(rawResponse), orderResponse)
	if dErr != nil {
		p.log.Info("Decoding error", dErr)

		return nil, errors.New("error decoding the order")
	}

	var redirectUrl string

	for _, url := range orderResponse.Links {
		found := url.Rel == "approve"
		if found {
			redirectUrl = url.Href
		}
	}

	return &PaymentDetail{
		PrivateId:   orderResponse.Id,
		RedirectUrl: redirectUrl,
		Status:      "CREATED",
	}, nil
}

func (p *PayPal) Capture(id string) (bool, error) {
	request, err := http.NewRequest(http.MethodPost, p.basePath+"/v2/checkout/orders/"+id+"/capture", nil)
	if err != nil {
		p.log.Error("Error on request", err)
	}

	response, err := p.requestWrapper(*request)
	if err != nil {
		p.log.Error("Do request err", err)
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		p.log.Error("Error decoding order response", rawResponse)
		return false, errors.New("error decoding order response")
	}

	defer response.Body.Close()

	isCreatedStatus := response.StatusCode == http.StatusCreated
	if !isCreatedStatus {
		p.log.Error("Error capturing the order", "response", string(rawResponse), "status", response.StatusCode)

		return false, errors.New("error capturing")
	}

	return true, nil
}

func (p *PayPal) Refund(paymentId string, refund PartialRefund) (*RefundResponse, error) {
	orderDetail, err := p.getOrder(paymentId)
	if err != nil {
		p.log.Warn("Order querying", "orderId", paymentId)
		return nil, errors.New("Error querying the order" + paymentId)
	}

	purchaseUnits := orderDetail.PurchaseUnits
	if purchaseUnits == nil {
		return nil, errors.New("order not have purchase units")
	}

	captures := purchaseUnits[0].Payments.Captures
	if captures == nil {
		return nil, errors.New("order not capture yet")
	}

	var payload Refund
	payload.Amount.Value = strconv.Itoa(int(refund.Amount / 100))
	payload.Amount.CurrencyCode = captures[0].Amount.CurrencyCode

	jsonMarshal, err := json.Marshal(payload)
	if err != nil {
		p.log.Error("error on marshal refund request")
	}

	request, err := http.NewRequest(http.MethodPost, p.basePath+"/v2/payments/captures/"+captures[0].ID+"/refund", bytes.NewBuffer(jsonMarshal))
	if err != nil {
		p.log.Error("error creating request for refund, " + paymentId)
	}

	response, err := p.requestWrapper(*request)
	if err != nil {
		p.log.Error("error requesting refund " + paymentId)
	}

	defer response.Body.Close()
	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		p.log.Error("Error reading body for refund: " + paymentId)

		return nil, errors.New("error refunding order with status " + request.Response.Status)
	}

	isOk := response.StatusCode == http.StatusCreated
	if !isOk {
		p.log.Error("error refunding order with status "+response.Status, "response", rawResponse)
		return nil, errors.New("error refunding order with status " + response.Status)
	}

	var refundDetail RefundDetail
	unmarshalErr := json.Unmarshal(rawResponse, &refundDetail)

	if unmarshalErr != nil {
		return nil, errors.New("error decoding json")
	}

	return &RefundResponse{
		Id:     refundDetail.Id,
		Amount: captures[0].Amount.CurrencyCode,
	}, nil
}

func (p *PayPal) getOrder(orderId string) (*OrderDetail, error) {
	request, err := http.NewRequest(http.MethodGet, p.basePath+"/v2/checkout/orders/"+orderId, nil)
	if err != nil {
		p.log.Error("error creating request for refund, " + orderId)
	}

	response, err := p.requestWrapper(*request)
	if err != nil {
		p.log.Error("error requesting refund " + orderId)
	}

	defer response.Body.Close()

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		p.log.Error("Error reading body for refund: " + orderId)

		return nil, errors.New("error refunding order with status " + response.Status)
	}

	var orderDetail OrderDetail

	marshalErr := json.Unmarshal(rawResponse, &orderDetail)
	if marshalErr != nil {
		return nil, errors.New("error decoding json")
	}

	return &orderDetail, nil

}

func (p *PayPal) getToken() (*string, error) {
	payload := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest(http.MethodPost, p.basePath+"/v1/oauth2/token", payload)
	if err != nil {
		p.log.Error("error creating the request", "detail", err)
		return nil, errors.New("error creating the request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.username, p.password)

	response, err := p.Client.Do(req)
	if err != nil {
		p.log.Error("error on request", "detail", err)
		return nil, errors.New("error on request")
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		p.log.Error("error on decoding", "detail", err)
		return nil, errors.New("error on decoding")
	}

	isOk := response.StatusCode == 200
	if !isOk {
		p.log.Error("Error getting the auth token", "detail", rawResponse)
		return nil, errors.New("error getting the auth token")
	}

	var tokenResponse TokenResponse
	unmarshalErr := json.Unmarshal(rawResponse, &tokenResponse)
	if unmarshalErr != nil {
		p.log.Error("error unmarshal", "detail", unmarshalErr)
	}

	return &tokenResponse.AccessToken, nil
}

func (p *PayPal) retryRequest(request http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", "Bearer "+p.bearerToken)

	response, err := p.Client.Do(&request)
	if err != nil {
		p.log.Error("RETRY_REQUEST", "message", err)

		return nil, errors.New("RETRY_REQUEST fails")
	}

	return response, nil
}

func (p *PayPal) requestWrapper(request http.Request) (*http.Response, error) {
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+p.bearerToken)

	var body []byte
	if request.Body != nil {
		body, _ = io.ReadAll(request.Body)
		defer request.Body.Close()
	}
	bodyCopy := io.NopCloser(bytes.NewReader(body))

	request.Body = bodyCopy
	firstResponse, err := p.Client.Do(&request)
	if err != nil {
		log.Println("Error requesting the token")
	}

	isUnauthorized := firstResponse.StatusCode == 401
	if isUnauthorized {
		token, err := p.getToken()

		if err != nil {
			return nil, errors.New("error getting authorization bearer token")
		}

		p.bearerToken = *token
		bodyCopy := io.NopCloser(bytes.NewReader(body))
		request.Body = bodyCopy

		return p.retryRequest(request)
	}

	return firstResponse, nil
}
