package main

import (
	"errors"
	"fmt"
	"time"
)

type InMemory struct {
}

var payments []Payment

func (InMemory) save(payment Payment) string {
	paymentId := fmt.Sprint(time.Now().UnixNano())
	payment.Id = paymentId
	payments = append(payments, payment)

	return paymentId
}

func (im InMemory) findById(id string) (*Payment, error) {
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

func (im InMemory) updateStatus(id string, status string) error {
	for index, p := range payments {
		match := id == p.Id
		if match {
			payments[index].Status = status
		}
	}

	return nil
}

func (im InMemory) attachRefund(paymentId string, refund RefundResponse) error {
	for index, p := range payments {
		match := paymentId == p.Id
		if match {
			payments[index].Refunds = append(payments[index].Refunds, refund)
		}
	}

	return nil
}
