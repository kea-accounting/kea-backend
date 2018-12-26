package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"

	"github.com/kea-accounting/kea-backend/data"
	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/kea-accounting/kea-backend/errors"
	"github.com/kea-accounting/kea-backend/globals"
	ohttp "github.com/kea-accounting/kea-backend/http"
	"github.com/kea-accounting/kea-backend/util"
)

func ListInvoices(w http.ResponseWriter, r *http.Request) {
	status := mux.Vars(r)["status"]
	userID := r.Context().Value(globals.UserIDKey).(string)
	invoices, err := data.GetInvoicesByStatus(userID, status)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {
		ohttp.WriteJSON(w, http.StatusOK, invoices)
	}
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(globals.UserIDKey).(string)
	// get id
	clientID := mux.Vars(r)["id"]
	if !util.IsValidName(clientID) {
		err := errors.BadRequest(
			fmt.Errorf("Invoice id '%s' is not a valid id", clientID))

		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	invoice, err := data.GetInvoice(clientID, userID)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {
		ohttp.WriteJSON(w, http.StatusOK, invoice)
	}
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(globals.UserIDKey).(string)

	var invoice datastructures.Invoice
	jsonParser := json.NewDecoder(r.Body)
	err := jsonParser.Decode(&invoice)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	invoice.Created, err = ptypes.TimestampProto(time.Now())
	invoice.LastUpdated, err = ptypes.TimestampProto(time.Now())
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	invoice.Id = util.NewID()

	_, err = data.SaveInvoice(&invoice, userID)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {
		ohttp.WriteJSON(w, http.StatusOK, invoice)
	}
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(globals.UserIDKey).(string)

	clientID := mux.Vars(r)["id"]
	var invoice datastructures.Invoice
	jsonParser := json.NewDecoder(r.Body)
	err := jsonParser.Decode(&invoice)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	existingInvoice, err := data.GetInvoice(clientID, userID)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	newInvoice, err := data.GetInvoice(clientID, userID)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	if existingInvoice.GetId() != invoice.GetId() {
		err := errors.BadRequest(
			fmt.Errorf("Invalid invoice id %s", invoice.GetId()))

		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	newInvoice.LastUpdated, err = ptypes.TimestampProto(time.Now())
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	newInvoice.ClientId = invoice.GetClientId()
	newInvoice.Lineitems = invoice.GetLineitems()
	newInvoice.Project = invoice.GetProject()
	newInvoice.Status = invoice.GetStatus()
	newInvoice.PaymentTerms = invoice.GetPaymentTerms()
	newInvoice.VatPeriod = invoice.GetVatPeriod()

	_, err = data.UpdateInvoice(newInvoice, existingInvoice)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {
		ohttp.WriteJSON(w, http.StatusOK, invoice)
	}
}

func TotalInvoicesByPeriod(w http.ResponseWriter, r *http.Request) {
	period := mux.Vars(r)["periodKey"]
	userID := r.Context().Value(globals.UserIDKey).(string)
	user, err := data.GetUserById(userID)

	totals, err := data.TotalInvoicesByPeriod(userID, period, user.GetVatFlatRate())
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {
		ohttp.WriteJSON(w, http.StatusOK, totals)
	}
}
