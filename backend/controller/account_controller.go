package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// AccountController : Account controller struct
type AccountController struct {
	Name string
}

// Create : create Account object for specified account
func (ctrl AccountController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	_, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating account: %v", err))
		return
	}
	account := &model.Account{}
	if err := utility.ReadRequestData(r, account); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating account: %v", err))
		return
	}
	if err := account.Insert(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating account: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(account))
}

// Find : find Account object by its ID
func (ctrl AccountController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	_, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding account: %v", err))
		return
	}
	accountID := pat.Param(r, "id")
	account := &model.Account{ID: accountID}
	if err := account.FindByID(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding account: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(account))
}

// FindAll : get all Account object
func (ctrl AccountController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all accounts: %v", err))
		return
	}
	accounts, err := model.AllAccounts(session.TenantID)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all accounts: %v", err))
		return
	}
	response := Response{Data: accounts}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update Account object
func (ctrl AccountController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	_, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating account: %v", err))
		return
	}
	accountID := pat.Param(r, "id")
	account := &model.Account{ID: accountID}
	if err := utility.ReadRequestData(r, account); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating account: %v", err))
		return
	}
	if err := account.UpdateAll(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating account: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(account))
}

// Delete : delete Account object by its ID
func (ctrl AccountController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	_, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating account: %v", err))
		return
	}
	accountID := pat.Param(r, "id")
	account := &model.Account{ID: accountID}
	if err := account.Delete(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting account: %v", err))
		return
	}
	respondSucceed(w, "Deleting Account successfully")
}
