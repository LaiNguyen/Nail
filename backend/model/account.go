package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

const (
	Worker    = 0
	Supervior = 1
	Admin     = 2
)

const (
	Unconfirmed = 0
	Confirmed   = 1
	Disabled    = 2
)

// Account struct
type Account struct {
	ID       string `json:"id" bson:"_id"`
	TenantID string `json:"tenant_id" bson:"tenant_id"`
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Role     int    `json:"role" bson:"role"`
	Status   int    `json:"status" bson:"status"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// AccountCollectionName returns the full name of Account collection
func AccountCollectionName() string {
	return fmt.Sprintf("account")
}

// Insert : insert new Account record
func (account *Account) Insert() error {
	//TODO: check the username is unique
	account.ID = bson.NewObjectId().Hex()
	account.UpdatedBy = account.CreatedBy
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Insert(account)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting Account", err)
	}

	return nil
}

// Update : udpate Account based on selector and updator
func (account *Account) Update(selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find Account by its ID
func (account *Account) FindByID() error {
	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": account.ID,
			}
			return collection.Find(selector).One(account)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting Account by id[%s]", err, account.ID)
	}

	return nil
}

// AllAccounts : find all accounts
func AllAccounts(tenantID string) ([]Account, error) {
	var accounts []Account
	selector := bson.M{
		"tenant_id": tenantID,
	}
	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Find(selector).All(&accounts)
		}); err != nil {
		return accounts, fmt.Errorf("Error[%s] while getting all accounts", err)
	}
	return accounts, nil
}

// Delete : delete Account by its ID
func (account *Account) Delete() error {
	if err := mongo.Execute("monotonic", AccountCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": account.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting Account by id[%s]", err, account.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (account *Account) UpdateAll() error {
	selector := bson.M{
		"_id": account.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"tenant_id": account.TenantID,
			"user_name": account.UserName,
			"password":  account.Password,
			"phone":     account.Phone,
			"role":      account.Role,
			"u_by":      account.UpdatedBy,
			"u_at":      time.Now(),
		},
	}
	return account.Update(selector, updator)
}
