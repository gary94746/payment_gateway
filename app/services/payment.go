package services

import (
	"errors"

	"payment-processor.gary94746/main/lib/database"
	"payment-processor.gary94746/main/lib/processors"
)

func (s *Services) CreatePayment(payment processors.Payment) (*processors.PaymentDetail, error) {
	paymentCreation, err := s.PaymentProcessor.Create(payment)
	if err != nil {
		return nil, errors.New("error creating the payment")
	}

	items := []database.LineItem{}
	for _, item := range payment.LineItems {
		items = append(items, database.LineItem{
			Name:     item.Name,
			Amount:   item.Amount,
			Quantity: item.Quantity,
		})
	}

	databasePayment := database.Payment{
		Currency:    payment.Currency,
		Amount:      payment.Amount,
		Status:      processors.StatusCreated,
		RedirectUrl: payment.RedirectUrl,
		CancelUrl:   payment.CancelUrl,
		PrivateId:   paymentCreation.PrivateId,
		Id:          payment.Id,
		LineItems:   items,
		Refunds:     []database.RefundResponse{},
	}

	paymentId := s.Database.Save(databasePayment)
	payment.Id = paymentId
	paymentCreation.Id = paymentId

	return paymentCreation, nil
}

func (s *Services) CapturePayment(paymentId string) error {
	payment, err := s.Database.FindById(paymentId)
	if err != nil {
		return errors.New("payment not found")
	}

	_, captureErr := s.PaymentProcessor.Capture(payment.PrivateId)
	if captureErr != nil {
		return errors.New(captureErr.Error())
	}

	s.Database.UpdateStatus(paymentId, processors.StatusCaptured)

	return nil
}

func (s *Services) GetPayment(paymentId string) (*database.Payment, error) {
	payment, err := s.Database.FindById(paymentId)

	if err != nil {
		return nil, errors.New("error getting payment")
	}

	return payment, nil
}

func (s *Services) RefundPayment(paymentId string, refund processors.PartialRefund) (*processors.RefundResponse, error) {
	order, err := s.Database.FindById(paymentId)
	if err != nil {
		return nil, err
	}

	refundRes, err1 := s.PaymentProcessor.Refund(order.PrivateId, refund)
	if err1 != nil {
		return nil, err1
	}

	s.Database.UpdateStatus(paymentId, processors.StatusRefunded)
	s.Database.AttachRefund(paymentId, database.RefundResponse{
		Id:     refundRes.Id,
		Amount: refundRes.Amount,
	})

	return refundRes, nil
}
