package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
	tools "nail/backend/tool"
)

// Session struct
type Session struct {
	ID          string    `json:"id" bson:"_id"`
	Token       string    `json:"token" bson:"token"`
	SessionType string    `json:"session_type" bson:"session_type"`
	TenantID    string    `json:"tenant_id" bson:"tenant_id"`
	AccountID   string    `json:"account_id" bson:"account_id"`
	ExpiredAt   time.Time `json:"e_at" bson:"e_at"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`

	// Account Account `json:"account" bson:"-"`
}

// SessionCollectionName returns the full name of session collection
func SessionCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_session", tenantID)
}

// Insert : insert new session record
func (session *Session) Insert(tenantID string) error {
	session.ID = bson.NewObjectId().Hex()
	session.Token = tools.UUID()
	session.ExpiredAt = time.Now().Add(time.Duration(1) * time.Hour)

	session.UpdatedBy = session.CreatedBy
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", SessionCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(session)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting session", err)
	}

	return nil
}

// Update : udpate session based on selector and updator
func (session *Session) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", SessionCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find session by its ID
func (session *Session) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", SessionCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": session.ID,
			}
			return collection.Find(selector).One(session)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting session by id[%s]", err, session.ID)
	}

	return nil
}

// AllSessions : find all sessions
func AllSessions(tenantID string) ([]Session, error) {
	var sessions []Session
	if err := mongo.Execute("monotonic", SessionCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Find(nil).All(&sessions)
		}); err != nil {
		return sessions, fmt.Errorf("Error[%s] while getting all sessions", err)
	}
	return sessions, nil
}

// Delete : delete session by its ID
func (session *Session) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", SessionCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": session.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting session by id[%s]", err, session.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (session *Session) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": session.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"token":        session.Token,
			"session_type": session.SessionType,
			"tenant_id":    session.TenantID,
			"account_id":   session.AccountID,
			"e_at":         session.ExpiredAt,
			"u_at":         time.Now(),
		},
	}
	return session.Update(tenantID, selector, updator)
}

// Login: process login request
func Login(userName, password string) (Session, error) {
	// TODO: add check login attempts count

	session := Session{}
	account := Account{UserName: userName, Password: password}

	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"user_name": account.UserName,
				"password":  account.Password,
			}
			return collection.Find(selector).One(&account)
		}); err != nil {
		return session, fmt.Errorf("Error[%s] while login with username[%s] and password[%s]", err, account.UserName, account.Password)
	}

	session.TenantID = account.TenantID
	session.AccountID = account.ID
	session.SessionType = "login"
	session.CreatedBy = account.Email
	if err := session.Insert(account.TenantID); err != nil {
		return session, fmt.Errorf("Error[%s] while create session with username[%s] and password[%s]", err, account.UserName, account.Password)
	}
	return session, nil
}

// Logout : kill the session
func (session *Session) Logout() error {
	return session.Delete(session.TenantID)
}

// IsValid : check if session is valid or not
func (session *Session) IsValid() error {
	if err := mongo.Execute("monotonic", SessionCollectionName(session.TenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"token": session.Token,
			}
			return collection.Find(selector).One(session)
		}); err != nil {
		return fmt.Errorf("Error[%s] while finding session with token[%s]", err, session.Token)
	}
	// if session.Expired() {
	// 	return fmt.Errorf("Session expired")
	// }
	// extend expired Time
	selector := bson.M{
		"_id": session.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"e_at": time.Now().Add(time.Duration(1) * time.Hour),
			"u_at": time.Now(),
		},
	}
	return session.Update(session.TenantID, selector, updator)
}

// Expired returns true if expired
func (session *Session) Expired() bool {
	return time.Now().After(session.ExpiredAt)
}
