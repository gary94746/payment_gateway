package processors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Stripe struct {
	Client   *http.Client
	token    string
	log      slog.Logger
	basePath string
}

func (s *Stripe) doRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Authorization", "Bearer "+s.token)

	response, err := s.Client.Do(request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Stripe) Init(settings PaymentSettings) error {
	s.log = *slog.New(slog.NewJSONHandler(os.Stdout, nil))
	s.basePath = "https://api.stripe.com/v1"
	s.token = settings.Credentials["token"]

	s.Client = &http.Client{
		Timeout: 60 * time.Second,
	}

	return nil
}

func (s *Stripe) Create(payment Payment) (*PaymentDetail, error) {
	form := url.Values{}
	for index, item := range payment.LineItems {
		form.Add(fmt.Sprintf("line_items[%d][amount]", index), strconv.Itoa(int(item.Amount)))
		form.Add(fmt.Sprintf("line_items[%d][currency]", index), payment.Currency)
		form.Add(fmt.Sprintf("line_items[%d][name]", index), item.Name)
		form.Add(fmt.Sprintf("line_items[%d][quantity]", index), strconv.Itoa(int(item.Quantity)))
	}

	form.Add("cancel_url", payment.CancelUrl)
	form.Add("success_url", payment.RedirectUrl)
	form.Add("mode", "payment")

	request, err := http.NewRequest(http.MethodPost, s.basePath+"/checkout/sessions", bytes.NewBuffer([]byte(form.Encode())))

	if err != nil {
		s.log.Error("error creating the request", "err", err.Error())
		return nil, errors.New("error creating the request")
	}

	response, err := s.doRequest(request)
	if err != nil {
		s.log.Error("error on request", "status", response.StatusCode)
		return nil, errors.New("error on request: " + err.Error())
	}
	defer response.Body.Close()

	decoded, err := io.ReadAll(response.Body)
	if err != nil {
		s.log.Error("error reading body", "err", err.Error())

		return nil, errors.New("error reading body")
	}

	isOk := response.StatusCode == http.StatusOK
	if !isOk {
		s.log.Warn("request fails", "status", response.StatusCode, "body", string(decoded))
		return nil, errors.New("Error requesting " + response.Status)
	}

	var checkout CheckoutResponse

	error := json.Unmarshal(decoded, &checkout)
	if error != nil {
		s.log.Error("error decoding json", "err", error.Error())
		return nil, errors.New("error decoding to json ")
	}

	return &PaymentDetail{
		PrivateId:   checkout.PaymentIntent,
		RedirectUrl: checkout.Url,
	}, nil
}

func (s *Stripe) Capture(id string) (bool, error) {
	intent, err := s.getPaymentIntent(id)
	if err != nil {
		return false, errors.New("error getting sessionId")
	}

	isPaid := intent.Status == "succeeded"
	if !isPaid {
		return false, errors.New("payment intent is not paid")
	}

	return true, nil
}

func (s *Stripe) Refund(paymentId string, refund PartialRefund) (*RefundResponse, error) {
	form := url.Values{}
	form.Add("payment_intent", paymentId)
	form.Add("amount", strconv.Itoa(int(refund.Amount)))

	request, err := http.NewRequest(http.MethodPost, s.basePath+"/refunds", bytes.NewBuffer([]byte(form.Encode())))
	if err != nil {
		return nil, errors.New("error creating the request")
	}

	response, err := s.doRequest(request)
	if err != nil {
		return nil, errors.New("error requesting refund")
	}

	rawPayload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading the payload")
	}

	defer response.Body.Close()

	isOk := response.StatusCode == http.StatusOK
	if !isOk {
		return nil, errors.New("response with status " + response.Status)
	}

	type CustomRefundResponse struct {
		Id     string `json:"id"`
		Amount int64  `json:"amount"`
	}

	var refundResponse CustomRefundResponse
	unmarshalError := json.Unmarshal(rawPayload, &refundResponse)
	if unmarshalError != nil {
		s.log.Error("decoding raw payload error", "message", string(rawPayload))
		return nil, errors.New("error parsing the response: " + unmarshalError.Error())
	}

	return &RefundResponse{
		Id:     refundResponse.Id,
		Amount: strconv.Itoa(int(refundResponse.Amount)),
	}, nil
}

func (s *Stripe) getPaymentIntent(intentId string) (*PaymentIntentResponse, error) {
	request, err := http.NewRequest(http.MethodGet, s.basePath+"/payment_intents/"+intentId, nil)
	if err != nil {
		return nil, errors.New("error creating request")
	}

	response, err := s.doRequest(request)
	if err != nil {
		return nil, errors.New("Error requesting" + response.Status)
	}

	isOk := response.StatusCode == http.StatusOK
	if !isOk {
		return nil, errors.New("error getting session information " + response.Status)
	}

	rawPayload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error decoding body")
	}

	defer response.Body.Close()

	var sessionDetail PaymentIntentResponse
	unmarshalError := json.Unmarshal(rawPayload, &sessionDetail)
	if unmarshalError != nil {
		s.log.Error("decoding error", "message", string(rawPayload))
		return nil, errors.New("error parsing to json")
	}

	return &sessionDetail, nil
}
