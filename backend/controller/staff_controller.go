package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// StaffController : staff controller struct
type StaffController struct {
	Name string
}

// Create : create staff object for specified account
func (ctrl StaffController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating staff: %v", err))
		return
	}
	staff := &model.Staff{}
	if err := utility.ReadRequestData(r, staff); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating staff: %v", err))
		return
	}
	if err := staff.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating staff: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(staff))
}

// Find : find staff object by its ID
func (ctrl StaffController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding staff: %v", err))
		return
	}
	staffID := pat.Param(r, "id")
	staff := &model.Staff{ID: staffID}
	if err := staff.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding staff: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(staff))
}

// FindAll : get all staff object
func (ctrl StaffController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all staffs: %v", err))
		return
	}
	filterInfo := model.FilterInfo{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo = GetFilterInfo(r)
	}
	staffs, total, err := model.AllStaffs(session.TenantID, filterInfo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all staffs: %v", err))
		return
	}
	response := Response{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo.Total = total
		response = PrepareResponse(filterInfo, staffs)
	} else {
		response.Data = staffs
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update staff object
func (ctrl StaffController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating staff: %v", err))
		return
	}
	staffID := pat.Param(r, "id")
	staff := &model.Staff{ID: staffID}
	if err := utility.ReadRequestData(r, staff); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating staff: %v", err))
		return
	}
	if err := staff.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating staff: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(staff))
}

// Delete : delete staff object by its ID
func (ctrl StaffController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting staff: %v", err))
		return
	}
	staffID := pat.Param(r, "id")
	staff := &model.Staff{ID: staffID}
	if err := staff.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting staff: %v", err))
		return
	}
	respondSucceed(w, "Deleting staff successfully")
}
