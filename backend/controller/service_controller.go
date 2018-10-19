package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// ServiceController : service controller struct
type ServiceController struct {
	Name string
}

// Create : create service object for specified account
func (ctrl ServiceController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating service: %v", err))
		return
	}

	service := &model.Service{}
	if err := utility.ReadRequestData(r, service); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating service: %v", err))
		return
	}

	if err := service.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating service: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(service))
}

// Find : find service object by its ID
func (ctrl ServiceController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding service: %v", err))
		return
	}
	serviceID := pat.Param(r, "id")
	service := &model.Service{ID: serviceID}
	if err := service.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding service: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(service))
}

// FindAll : get all service object
func (ctrl ServiceController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding service: %v", err))
		return
	}
	filterInfo := model.FilterInfo{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo = GetFilterInfo(r)
	}
	services, total, err := model.AllServices(session.TenantID, filterInfo)
	response := Response{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo.Total = total
		response = PrepareResponse(filterInfo, services)
	} else {
		response.Data = services
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update service object
func (ctrl ServiceController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating service: %v", err))
		return
	}
	serviceID := pat.Param(r, "id")
	service := &model.Service{ID: serviceID}
	if err := utility.ReadRequestData(r, service); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating service: %v", err))
		return
	}
	if err := service.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating service: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(service))
}

// Delete : delete service object by its ID
func (ctrl ServiceController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting service: %v", err))
		return
	}
	serviceID := pat.Param(r, "id")
	service := &model.Service{ID: serviceID}
	if err := service.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting service: %v", err))
		return
	}
	respondSucceed(w, "Deleting service successfully")
}
