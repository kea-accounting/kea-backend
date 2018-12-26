package data

import (
	"fmt"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/kea-accounting/kea-backend/domain"
	"github.com/kea-accounting/kea-backend/errors"
	"github.com/shopspring/decimal"
)

type Totals struct {
	VatDueSales                  decimal.Decimal `json:"vatDueSales"`
	VatDueAcquisitions           decimal.Decimal `json:"vatDueAcquisitions"`
	TotalVatDue                  decimal.Decimal `json:"totalVatDue"`
	VatReclaimedCurrPeriod       decimal.Decimal `json:"vatReclaimedCurrPeriod"`
	NetVatDuenetVatDue           decimal.Decimal `json:"netVatDue"`
	TotalValueSalesExVAT         decimal.Decimal `json:"totalValueSalesExVAT"`
	TotalValuePurchasesExVAT     decimal.Decimal `json:"totalValuePurchasesExVAT"`
	TotalValueGoodsSuppliedExVAT decimal.Decimal `json:"totalValueGoodsSuppliedExVAT"`
	TotalAcquisitionsExVAT       decimal.Decimal `json:"totalAcquisitionsExVAT"`
}

func SaveInvoice(invoice *datastructures.Invoice, userID string) (*datastructures.Invoice, error) {

	_, err := domain.Totalize(invoice)
	if err != nil {
		return nil, errors.WrapError(err)
	}
	invoice.UserId = userID

	i, err := RunTransaction(invoice, func(out *[]byte, tr *fdb.Transaction) (interface{}, error) {

		existing := tr.Get(subspaces.invoiceByClientID.Pack(tuple.Tuple{invoice.GetUserId(), invoice.GetClientId()})).MustGet()

		if existing != nil {
			return nil, errors.Conflict(
				fmt.Errorf("Invoice with ID '%s' already exists", invoice.GetClientId()))
		}

		tr.Set(subspaces.invoiceByClientID.Pack(tuple.Tuple{invoice.GetUserId(), invoice.GetClientId()}), subspaces.invoice.Pack(tuple.Tuple{invoice.GetId()}))
		tr.Set(subspaces.invoiceByStatus.Pack(tuple.Tuple{invoice.GetUserId(), invoice.GetStatus().String(), invoice.GetDateDue(), invoice.GetId()}), []byte{})
		tr.Set(subspaces.invoiceByVATPeriod.Pack(tuple.Tuple{invoice.GetUserId(), invoice.GetVatPeriod(), invoice.GetDateDue(), invoice.GetId()}), []byte{})
		tr.Set(subspaces.invoice.Pack(tuple.Tuple{invoice.GetId()}), *out)

		return invoice, nil
	})

	if err != nil {
		return nil, err
	}

	return i.(*datastructures.Invoice), nil
}

func UpdateInvoice(newInvoice, existingInvoice *datastructures.Invoice) (*datastructures.Invoice, error) {

	_, err := domain.Totalize(existingInvoice)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	i, err := RunTransaction(newInvoice, func(out *[]byte, tr *fdb.Transaction) (interface{}, error) {
		tr.Clear(subspaces.invoiceByClientID.Pack(tuple.Tuple{existingInvoice.GetUserId(), existingInvoice.GetClientId()}))
		tr.Set(subspaces.invoiceByClientID.Pack(tuple.Tuple{newInvoice.GetUserId(), newInvoice.GetClientId()}), subspaces.invoice.Pack(tuple.Tuple{newInvoice.GetId()}))

		tr.Clear(subspaces.invoiceByStatus.Pack(tuple.Tuple{existingInvoice.GetUserId(), existingInvoice.GetStatus().String(), existingInvoice.GetDateDue(), newInvoice.GetId()}))
		tr.Set(subspaces.invoiceByStatus.Pack(tuple.Tuple{newInvoice.GetUserId(), newInvoice.GetStatus().String(), newInvoice.GetDateDue(), newInvoice.GetId()}), []byte{})

		tr.Clear(subspaces.invoiceByVATPeriod.Pack(tuple.Tuple{existingInvoice.GetUserId(), existingInvoice.GetVatPeriod(), existingInvoice.GetDateDue(), newInvoice.GetId()}))
		tr.Set(subspaces.invoiceByVATPeriod.Pack(tuple.Tuple{newInvoice.GetUserId(), newInvoice.GetVatPeriod(), newInvoice.GetDateDue(), newInvoice.GetId()}), []byte{})

		tr.Set(subspaces.invoice.Pack(tuple.Tuple{newInvoice.GetId()}), *out)

		return newInvoice, nil
	})

	if err != nil {
		return nil, err
	}

	return i.(*datastructures.Invoice), nil
}

func GetInvoice(clientID string, userID string) (*datastructures.Invoice, error) {
	bytes, err := GetByLookup(subspaces.invoiceByClientID.Pack(tuple.Tuple{userID, clientID}))
	if err != nil {
		return nil, errors.WrapError(err)
	}
	if bytes == nil {
		return nil, errors.NotFound(
			fmt.Errorf("Invoice with ID '%s' does not exist", clientID))
	}

	invoice := &datastructures.Invoice{}
	if err := proto.Unmarshal(bytes, invoice); err != nil {
		return nil, errors.WrapError(err)
	}

	return invoice, nil
}

func GetInvoicesByStatus(userID, status string) ([]*datastructures.Invoice, error) {
	i, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {

		var invoices []*datastructures.Invoice
		ri := tr.GetRange(subspaces.invoiceByStatus.Sub(userID, status), fdb.RangeOptions{}).Iterator()
		for ri.Advance() {
			kv := ri.MustGet()
			t, err := subspaces.invoiceByStatus.Unpack(kv.Key)

			id := subspaces.invoice.Pack(tuple.Tuple{t[3].(string)})
			bytes, err := GetById(id)
			if err != nil {
				return nil, errors.WrapError(err)
			}

			invoice := &datastructures.Invoice{}
			if err := proto.Unmarshal(bytes, invoice); err != nil {
				return nil, errors.WrapError(err)
			}

			invoices = append(invoices, invoice)
		}

		return invoices, nil
	})

	if err != nil {
		return nil, err
	}

	return i.([]*datastructures.Invoice), err
}

func TotalInvoicesByPeriod(userID, period, vatFlatRate string) (*Totals, error) {

	i, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {

		totals := &Totals{
			VatDueSales:                  decimal.New(0, 0),
			VatDueAcquisitions:           decimal.New(0, 0),
			TotalVatDue:                  decimal.New(0, 0),
			VatReclaimedCurrPeriod:       decimal.New(0, 0),
			NetVatDuenetVatDue:           decimal.New(0, 0),
			TotalValueSalesExVAT:         decimal.New(0, 0),
			TotalValuePurchasesExVAT:     decimal.New(0, 0),
			TotalValueGoodsSuppliedExVAT: decimal.New(0, 0),
			TotalAcquisitionsExVAT:       decimal.New(0, 0),
		}
		ri := tr.GetRange(subspaces.invoiceByVATPeriod.Sub(userID, period), fdb.RangeOptions{}).Iterator()
		for ri.Advance() {
			kv := ri.MustGet()
			t, err := subspaces.invoiceByVATPeriod.Unpack(kv.Key)

			id := subspaces.invoice.Pack(tuple.Tuple{t[3].(string)})
			bytes, err := GetById(id)
			if err != nil {
				return nil, errors.WrapError(err)
			}

			invoice := &datastructures.Invoice{}
			if err := proto.Unmarshal(bytes, invoice); err != nil {
				return nil, errors.WrapError(err)
			}

			invoiceTotal, err := decimal.NewFromString(invoice.GetTotalAmount())
			if err != nil {
				return nil, errors.WrapError(err)
			}
			totals.TotalValueSalesExVAT = totals.TotalValueSalesExVAT.Add(invoiceTotal)
		}

		flatRate, err := decimal.NewFromString(vatFlatRate)
		if err != nil {
			return nil, errors.WrapError(err)
		}
		flatRate = flatRate.Div(decimal.New(100, 0))
		totals.VatDueSales = totals.TotalValueSalesExVAT.Mul(flatRate)

		return totals, nil
	})

	if err != nil {
		return nil, err
	}

	return i.(*Totals), err
}
