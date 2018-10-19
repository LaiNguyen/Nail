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

// OrderController : order controller struct
type OrderController struct {
	Name string
}

// Create : create order object for specified account
func (ctrl OrderController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating order: %v", err))
		return
	}
	order := &model.Order{}
	if err := utility.ReadRequestData(r, order); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating order: %v", err))
		return
	}
	if err := order.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating order: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(order))
}

// Find : find order object by its ID
func (ctrl OrderController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding order: %v", err))
		return
	}
	orderID := pat.Param(r, "id")
	order := &model.Order{ID: orderID}
	if err := order.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding order: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(order))
}

// FindAll : get all order object
func (ctrl OrderController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding all orders: %v", err))
		return
	}
	queryValues := r.URL.Query()
	orderStatus, err := strconv.Atoi(queryValues.Get("status"))
	if err != nil {
		orderStatus = -1
	}
	orders, err := model.AllOrders(session.TenantID, orderStatus)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while fiding all orders: %v", err))
		return
	}
	response := Response{Data: orders}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update order object
func (ctrl OrderController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating orders: %v", err))
		return
	}
	orderID := pat.Param(r, "id")
	order := &model.Order{ID: orderID}
	if err := utility.ReadRequestData(r, order); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating order: %v", err))
		return
	}
	if err := order.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating order: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(order))
}

// Delete : delete order object by its ID
func (ctrl OrderController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting order: %v", err))
		return
	}
	orderID := pat.Param(r, "id")
	order := &model.Order{ID: orderID}
	if err := order.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting order: %v", err))
		return
	}
	respondSucceed(w, "Deleting order successfully")
}

type CheckoutData struct {
	Order   model.Order   `json:"order"`
	Payment model.Payment `json:"payment"`
}

// Checkout : checkout order object
func (ctrl OrderController) Checkout(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while checkout order: %v", err))
		return
	}
	checkoutData := &CheckoutData{}
	if err := utility.ReadRequestData(r, checkoutData); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while checkout order: %v", err))
		return
	}
	if err := checkoutData.Order.Checkout(session.TenantID, checkoutData.Payment); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while checkout order: %v", err))
		return
	}
	respondSucceed(w, "Checkout order successfully")
	// fmt.Fprintf(w, "%s", utility.ToJSON(checkoutData))
}
