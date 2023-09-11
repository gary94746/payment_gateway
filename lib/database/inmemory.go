package database

import (
	"errors"
	"fmt"
	"time"
)

type InMemory struct {
}

var payments []Payment

func (im InMemory) Save(payment Payment) string {
	paymentId := fmt.Sprint(time.Now().UnixNano())
	payment.Id = paymentId
	payments = append(payments, payment)

	return paymentId
}

func (im InMemory) FindById(id string) (*Payment, error) {
	var payment *Payment

	for _, p := range payments {
		match := id == p.Id
		if match {
			payment = &p
		}
	}

	if payment == nil {
		return nil, errors.New("payment not exists")
	}

	return payment, nil
}

func (im InMemory) UpdateStatus(id string, status string) error {
	for index, p := range payments {
		match := id == p.Id
		if match {
			payments[index].Status = status
		}
	}

	return nil
}

func (im InMemory) AttachRefund(paymentId string, refund RefundResponse) error {
	for index, p := range payments {
		match := paymentId == p.Id
		if match {
			payments[index].Refunds = append(payments[index].Refunds, refund)
		}
	}

	return nil
}
