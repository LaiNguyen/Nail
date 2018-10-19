package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// CustomerController : customer controller struct
type CustomerController struct {
	Name string
}

// Create : create customer object for specified account
func (ctrl CustomerController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating customer: %v", err))
		return
	}
	customer := &model.Customer{}
	if err := utility.ReadRequestData(r, customer); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating customer: %v", err))
		return
	}
	if err := customer.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating customer: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(customer))
}

// Find : find customer object by its ID
func (ctrl CustomerController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding customer: %v", err))
		return
	}
	customerID := pat.Param(r, "id")
	customer := &model.Customer{ID: customerID}
	if err := customer.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding customer: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(customer))
}

// FindAll : get all customer object
func (ctrl CustomerController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating all customers: %v", err))
		return
	}
	filterInfo := model.FilterInfo{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo = GetFilterInfo(r)
	}
	customers, total, err := model.AllCustomers(session.TenantID, filterInfo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating all customers: %v", err))
		return
	}
	response := Response{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo.Total = total
		response = PrepareResponse(filterInfo, customers)
	} else {
		response.Data = customers
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update customer object
func (ctrl CustomerController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating customer: %v", err))
		return
	}
	customerID := pat.Param(r, "id")
	customer := &model.Customer{ID: customerID}
	if err := utility.ReadRequestData(r, customer); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating customer: %v", err))
		return
	}
	if err := customer.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating customer: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(customer))
}

// Delete : delete customer object by its ID
func (ctrl CustomerController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting customer: %v", err))
		return
	}
	customerID := pat.Param(r, "id")
	customer := &model.Customer{ID: customerID}
	if err := customer.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting customer: %v", err))
		return
	}
	respondSucceed(w, "Deleting customer successfully")
}
