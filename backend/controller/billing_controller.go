package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// BillingController : billing controller struct
type BillingController struct {
	Name string
}

// Create : create billing object for specified account
func (ctrl BillingController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating billing: %v", err))
		return
	}
	billing := &model.Billing{}
	if err := utility.ReadRequestData(r, billing); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating billing: %v", err))
		return
	}
	if err := billing.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating billing: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(billing))
}

// Find : find billing object by its ID
func (ctrl BillingController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding billing: %v", err))
		return
	}
	billingID := pat.Param(r, "id")
	billing := &model.Billing{ID: billingID}
	if err := billing.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding billing: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(billing))
}

// FindAll : get all billing object
func (ctrl BillingController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while fiding all billings: %v", err))
		return
	}
	filterInfo := GetFilterInfo(r)
	billings, total, err := model.AllBillings(session.TenantID, filterInfo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while fiding all billings: %v", err))
		return
	}
	filterInfo.Total = total
	response := PrepareResponse(filterInfo, billings)

	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update billing object
func (ctrl BillingController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating billing: %v", err))
		return
	}
	billingID := pat.Param(r, "id")
	billing := &model.Billing{ID: billingID}
	if err := utility.ReadRequestData(r, billing); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating billing: %v", err))
		return
	}
	if err := billing.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating billing: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(billing))
}

// Delete : delete billing object by its ID
func (ctrl BillingController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting billing: %v", err))
		return
	}
	billingID := pat.Param(r, "id")
	billing := &model.Billing{ID: billingID}
	if err := billing.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting billing: %v", err))
		return
	}
	respondSucceed(w, "Deleting billing successfully")
}

func GetFilterInfo(r *http.Request) model.FilterInfo {
	queryValues := r.URL.Query()
	status, err := strconv.Atoi(queryValues.Get("status"))
	if err != nil {
		status = -1
	}
	sort := queryValues.Get("sort")
	filter := queryValues.Get("filter")
	page, _ := strconv.Atoi(queryValues.Get("page"))
	perPage, _ := strconv.Atoi(queryValues.Get("per_page"))
	filterInfo := model.FilterInfo{
		Status:  status,
		Sort:    sort,
		Filter:  filter,
		Page:    page,
		PerPage: perPage,
	}
	return filterInfo
}

func PrepareResponse(filterInfo model.FilterInfo, data interface{}) Response {
	var lastPage int
	if filterInfo.Total <= filterInfo.PerPage {
		lastPage = 1
	} else {
		lastPage = filterInfo.Total/filterInfo.PerPage + 1
	}

	response := Response{
		Total:       filterInfo.Total,
		PerPage:     filterInfo.PerPage,
		CurrentPage: filterInfo.Page,
		LastPage:    lastPage,
		NextPageURL: "",
		PrevPageURL: "",
		From:        (filterInfo.Page-1)*filterInfo.PerPage + 1,
		To:          filterInfo.Page * filterInfo.PerPage,
		Data:        data,
	}
	return response
}
