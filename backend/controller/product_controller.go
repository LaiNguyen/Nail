package controller

import (
	"fmt"
	"net/http"
	"time"

	"goji.io/pat"

	"nail/backend/model"
	"nail/backend/utility"
)

// ProductController : product controller struct
type ProductController struct {
	Name string
}

// Create : create product object for specified account
func (ctrl ProductController) Create(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating product: %v", err))
		return
	}
	product := &model.Product{}
	if err := utility.ReadRequestData(r, product); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating product: %v", err))
		return
	}
	if err := product.Insert(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while creating product: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(product))
}

// Find : find product object by its ID
func (ctrl ProductController) Find(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding product: %v", err))
		return
	}
	productID := pat.Param(r, "id")
	product := &model.Product{ID: productID}
	if err := product.FindByID(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding product: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(product))
}

// FindByService : find product object by service ID
func (ctrl ProductController) FindByService(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding products by service: %v", err))
		return
	}
	filterInfo := model.FilterInfo{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo = GetFilterInfo(r)
	}
	serviceID := pat.Param(r, "id")
	products, total, err := model.FindProductsByServiceID(session.TenantID, serviceID, filterInfo)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while finding products by service: %v", err))
		return
	}
	response := Response{}
	if r.URL.Query().Get("paging") == "yes" {
		filterInfo.Total = total
		response = PrepareResponse(filterInfo, products)
	} else {
		response.Data = products
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// FindAll : get all product object
func (ctrl ProductController) FindAll(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while fiding all products: %v", err))
		return
	}
	products, err := model.AllProducts(session.TenantID)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while fiding all products: %v", err))
		return
	}
	response := Response{Data: products}
	fmt.Fprintf(w, "%s", utility.ToJSON(response))
}

// Update : update product object
func (ctrl ProductController) Update(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating product: %v", err))
		return
	}
	productID := pat.Param(r, "id")
	product := &model.Product{ID: productID}
	if err := utility.ReadRequestData(r, product); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating product: %v", err))
		return
	}
	if err := product.UpdateAll(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while updating product: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(product))
}

// Delete : delete product object by its ID
func (ctrl ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting product: %v", err))
		return
	}
	productID := pat.Param(r, "id")
	product := &model.Product{ID: productID}
	if err := product.Delete(session.TenantID); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while deleting product: %v", err))
		return
	}
	respondSucceed(w, "Deleting product successfully")
}
