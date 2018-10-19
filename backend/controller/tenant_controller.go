package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// TenantController : tenant controller struct
type TenantController struct {
	Name string
}

// Create : create tenant object for specified account
func (ctrl TenantController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	tenant := &model.Tenant{}
	if err := utility.ReadRequestData(r, tenant); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating tenant: %v", err))
		return
	}
	if err := tenant.Insert(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating tenant: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(tenant))
}

// Find : find tenant object by its ID
func (ctrl TenantController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	tenantID := pat.Param(r, "id")
	tenant := &model.Tenant{ID: tenantID}
	err := tenant.FindByID()
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding tenant: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(tenant))
}

// FindAll : get all tenant object
func (ctrl TenantController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	tenants, err := model.AllTenants()
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all tenants: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(tenants))
}

// Update : update tenant object
func (ctrl TenantController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	tenantID := pat.Param(r, "id")
	tenant := &model.Tenant{ID: tenantID}
	if err := utility.ReadRequestData(r, tenant); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating tenant: %v", err))
		return
	}
	if err := tenant.UpdateAll(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating tenant: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(tenant))
}

// Delete : delete tenant object by its ID
func (ctrl TenantController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	tenantID := pat.Param(r, "id")
	tenant := &model.Tenant{ID: tenantID}
	if err := tenant.Delete(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting tenant: %v", err))
		return
	}
	respondSucceed(w, "Deleting tenant successfully")
}
