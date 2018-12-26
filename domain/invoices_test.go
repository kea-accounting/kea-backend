package domain

import (
	"testing"

	"github.com/kea-accounting/kea-backend/datastructures"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assert")
	}
}

func TestTotalize(t *testing.T) {
	invoice := &datastructures.Invoice{
		DateInvoice:  "2018-12-23",
		PaymentTerms: 28,
		Lineitems: []*datastructures.Invoice_LineItem{
			&datastructures.Invoice_LineItem{
				Description:  "Programming",
				Quantity:     2,
				Price:        "200.21",
				SalesTaxRate: 20,
			},
			&datastructures.Invoice_LineItem{
				Description:  "Consulting",
				Quantity:     5,
				Price:        "280.21",
				SalesTaxRate: 20,
			},
		},
	}
	Totalize(invoice)
	assert(t, invoice.GetCreated() != nil)
	assert(t, invoice.DateDue == "2019-01-20")
	assert(t, invoice.Status == datastructures.Invoice_NEW)
	assert(t, invoice.Lineitems[0].GetNetAmount() == "400.42")
	assert(t, invoice.Lineitems[0].GetTaxAmount() == "80.08")
	assert(t, invoice.Lineitems[0].GetTotalAmount() == "480.50")
	assert(t, invoice.Lineitems[1].GetNetAmount() == "1401.05")
	assert(t, invoice.Lineitems[1].GetTaxAmount() == "280.21")
	assert(t, invoice.Lineitems[1].GetTotalAmount() == "1681.26")
	assert(t, invoice.GetNetAmount() == "1801.47")
	assert(t, invoice.GetSalesTax() == "360.29")
	assert(t, invoice.GetTotalAmount() == "2161.76")
}
