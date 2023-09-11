package processors

import "time"

type Refund struct {
	Amount Amount `json:"amount"`
}

type RefundDetail struct {
	Id string `json:"id"`
}

type Amount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UnitAmount  Amount `json:"unit_amount"`
	Quantity    string `json:"quantity"`
}
type Breakdown struct {
	ItemTotal Amount `json:"item_total"`
}

type PurchaseUnitAmount struct {
	CurrencyCode string    `json:"currency_code"`
	Value        string    `json:"value"`
	Breakdown    Breakdown `json:"breakdown"`
}

type PurchaseUnits struct {
	Amount PurchaseUnitAmount `json:"amount"`
	Items  []Item             `json:"items"`
}

type ApplicationContext struct {
	ReturnUrl string `json:"return_url"`
	CancelUrl string `json:"cancel_url"`
}

type Order struct {
	Intent             string             `json:"intent"`
	ApplicationContext ApplicationContext `json:"application_context"`
	PurchaseUnits      []PurchaseUnits    `json:"purchase_units"`
}

type OrderResponse struct {
	Id     string              `json:"id"`
	Status string              `json:"status"`
	Links  []OrderResponseLink `json:"links"`
}

type OrderResponseLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OrderDetail struct {
	ID            string `json:"id"`
	Intent        string `json:"intent"`
	Status        string `json:"status"`
	PaymentSource struct {
		Paypal struct {
			EmailAddress  string `json:"email_address"`
			AccountID     string `json:"account_id"`
			AccountStatus string `json:"account_status"`
			Name          struct {
				GivenName string `json:"given_name"`
				Surname   string `json:"surname"`
			} `json:"name"`
			PhoneNumber struct {
				NationalNumber string `json:"national_number"`
			} `json:"phone_number"`
			Address struct {
				CountryCode string `json:"country_code"`
			} `json:"address"`
			Attributes struct {
				CobrandedCards []struct {
					Labels []interface{} `json:"labels"`
					Payee  struct {
						EmailAddress string `json:"email_address"`
						MerchantID   string `json:"merchant_id"`
					} `json:"payee"`
					Amount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"amount"`
				} `json:"cobranded_cards"`
			} `json:"attributes"`
		} `json:"paypal"`
	} `json:"payment_source"`
	PurchaseUnits []struct {
		ReferenceID string `json:"reference_id"`
		Amount      struct {
			CurrencyCode string `json:"currency_code"`
			Value        string `json:"value"`
			Breakdown    struct {
				ItemTotal struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"item_total"`
				Shipping struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"shipping"`
				Handling struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"handling"`
				Insurance struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"insurance"`
				ShippingDiscount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"shipping_discount"`
				Discount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"discount"`
			} `json:"breakdown"`
		} `json:"amount"`
		Payee struct {
			EmailAddress string `json:"email_address"`
			MerchantID   string `json:"merchant_id"`
		} `json:"payee"`
		Description    string `json:"description"`
		SoftDescriptor string `json:"soft_descriptor"`
		Items          []struct {
			Name       string `json:"name"`
			UnitAmount struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"unit_amount"`
			Tax struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"tax"`
			Quantity    string `json:"quantity"`
			Description string `json:"description"`
			Sku         string `json:"sku"`
		} `json:"items"`
		Shipping struct {
			Name struct {
				FullName string `json:"full_name"`
			} `json:"name"`
			Address struct {
				AddressLine1 string `json:"address_line_1"`
				AddressLine2 string `json:"address_line_2"`
				AdminArea2   string `json:"admin_area_2"`
				AdminArea1   string `json:"admin_area_1"`
				PostalCode   string `json:"postal_code"`
				CountryCode  string `json:"country_code"`
			} `json:"address"`
		} `json:"shipping"`
		Payments struct {
			Captures []struct {
				ID     string `json:"id"`
				Status string `json:"status"`
				Amount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"amount"`
				FinalCapture     bool `json:"final_capture"`
				SellerProtection struct {
					Status            string   `json:"status"`
					DisputeCategories []string `json:"dispute_categories"`
				} `json:"seller_protection"`
				SellerReceivableBreakdown struct {
					GrossAmount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"gross_amount"`
					PaypalFee struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"paypal_fee"`
					NetAmount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"net_amount"`
				} `json:"seller_receivable_breakdown"`
				Links []struct {
					Href   string `json:"href"`
					Rel    string `json:"rel"`
					Method string `json:"method"`
				} `json:"links"`
				CreateTime time.Time `json:"create_time"`
				UpdateTime time.Time `json:"update_time"`
			} `json:"captures"`
		} `json:"payments"`
	} `json:"purchase_units"`
	Payer struct {
		Name struct {
			GivenName string `json:"given_name"`
			Surname   string `json:"surname"`
		} `json:"name"`
		EmailAddress string `json:"email_address"`
		PayerID      string `json:"payer_id"`
		Phone        struct {
			PhoneNumber struct {
				NationalNumber string `json:"national_number"`
			} `json:"phone_number"`
		} `json:"phone"`
		Address struct {
			CountryCode string `json:"country_code"`
		} `json:"address"`
	} `json:"payer"`
	UpdateTime time.Time `json:"update_time"`
	Links      []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
	CreditFinancingOffer struct {
		Issuer             string `json:"issuer"`
		InstallmentDetails struct {
			PaymentDue struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"payment_due"`
		} `json:"installment_details"`
		Term int `json:"term"`
	} `json:"credit_financing_offer"`
}
