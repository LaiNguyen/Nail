package controller

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"nail/backend/model"
	"nail/backend/utility"
)

const HeaderAuthorization = "Authorization"

type LoginInfo struct {
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
}

// Response : Response struct
type Response struct {
	Total       int         `json:"total"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	NextPageURL string      `json:"next_page_url"`
	PrevPageURL string      `json:"prev_page_url"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	Data        interface{} `json:"data"`
}

// SessionController : service controller struct
type SessionController struct {
	Name string
}

// Login : login handler with username and password
func (ctrl SessionController) Login(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	loginInfo := &LoginInfo{}
	if err := utility.ReadRequestData(r, loginInfo); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while trying to login: %v", err))
		return
	}

	session, err := model.Login(loginInfo.UserName, loginInfo.Password)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while trying to login: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(session))
}

// Logout : logout/delete session
func (ctrl SessionController) Logout(w http.ResponseWriter, r *http.Request) {
	logStartInfo(ctrl.Name, funcName(), time.Now())

	session, err := Authenticate(r)
	if err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while trying to logout: %v", err))
		return
	}

	if err := session.Logout(); err != nil {
		responseError(w, http.StatusInternalServerError, fmt.Sprintf("Error while trying to logout: %v", err))
		return
	}
	fmt.Fprintf(w, "%s", utility.ToJSON(session))
}

// Authenticate: middleware function of controllers level, extract token info and return tenantID
func Authenticate(r *http.Request) (model.Session, error) {
	session := model.Session{}
	tokenInfo := r.Header.Get(HeaderAuthorization)
	if tokenInfo == "" {
		return session, fmt.Errorf("Token not provided")
	}
	tokenTenantID := strings.Split(tokenInfo, "|")
	if len(tokenTenantID) != 2 {
		return session, fmt.Errorf("Invalid token format")
	}
	session.Token = tokenTenantID[0]
	session.TenantID = tokenTenantID[1]
	if err := session.IsValid(); err != nil {
		return session, fmt.Errorf("Error [%s] while checking if session is valid", err)
	}
	return session, nil
}

func responseError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	fmt.Fprint(w, message)
}

func respondSucceed(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, "%s", utility.ToJSON(struct {
		Message string `json:"message"`
	}{
		Message: message,
	}))
}

func logStartInfo(ctrlName, funcName string, startTime time.Time) {
	log.Printf("%s execute %s at %s", ctrlName, funcName, startTime)
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
