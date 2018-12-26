package domain

import (
	"fmt"
	"time"

	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/shopspring/decimal"
)

func Totalize(invoice *datastructures.Invoice) (*datastructures.Invoice, error) {

	layout := "2006-01-02"
	t, err := time.Parse(layout, invoice.DateInvoice)
	if err != nil {
		return nil, err
	}
	duration := fmt.Sprintf("%dh", invoice.PaymentTerms*24)
	paymentDuration, err := time.ParseDuration(duration)
	dueDate := t.Add(paymentDuration)
	invoice.DateDue = dueDate.Format(layout)

	netAmount := decimal.New(0, 0)
	salesTax := decimal.New(0, 0)
	for _, item := range invoice.Lineitems {
		price, err := decimal.NewFromString(item.GetPrice())
		if err != nil {
			return nil, err
		}
		qty := decimal.New(item.GetQuantity(), 0)
		salesTaxRate := decimal.New(item.GetSalesTaxRate(), 0).Div(decimal.New(100, 0))

		itemAmount := price.Mul(qty)
		itemTaxAmount := itemAmount.Mul(salesTaxRate)
		netAmount = netAmount.Add(itemAmount)
		salesTax = salesTax.Add(itemTaxAmount)
		item.NetAmount = itemAmount.StringFixedBank(2)
		item.TaxAmount = itemTaxAmount.StringFixedBank(2)
		item.TotalAmount = itemAmount.Add(itemTaxAmount).StringFixedBank(2)
	}
	totalAmount := netAmount.Add(salesTax)

	invoice.NetAmount = netAmount.StringFixedBank(2)
	invoice.SalesTax = salesTax.StringFixedBank(2)
	invoice.TotalAmount = totalAmount.StringFixedBank(2)
	return invoice, nil
}
