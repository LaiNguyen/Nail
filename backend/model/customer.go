package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

const (
	New     = 0
	Normal  = 1
	Regular = 2
	VIP     = 3
)

// Customer struct
type Customer struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Gender   string `json:"gender" bson:"gender"`
	Birthday string `json:"birthday" bson:"birthday"`
	Status   int    `json:"status" bson:"status"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// CustomerCollectionName returns the full name of customer collection
func CustomerCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_customer", tenantID)
}

// Insert : insert new customer record
func (customer *Customer) Insert(tenantID string) error {
	customer.ID = bson.NewObjectId().Hex()
	customer.UpdatedBy = customer.CreatedBy
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", CustomerCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(customer)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting customer", err)
	}

	return nil
}

// Update : udpate customer based on selector and updator
func (customer *Customer) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", CustomerCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find customer by its ID
func (customer *Customer) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", CustomerCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": customer.ID,
			}
			return collection.Find(selector).One(customer)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting customer by id[%s]", err, customer.ID)
	}

	return nil
}

// AllCustomers : find all customers
func AllCustomers(tenantID string, filterInfo FilterInfo) ([]Customer, int, error) {
	total := 0
	var customers []Customer
	if err := mongo.Execute("monotonic", CustomerCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			total, _ = collection.Find(nil).Count()
			if (filterInfo == FilterInfo{}) {
				return collection.Find(nil).All(&customers)
			}
			sort := filterInfo.PrepareSort()
			return collection.Find(nil).
				Limit(filterInfo.PerPage).
				Skip((filterInfo.Page - 1) * filterInfo.PerPage).
				Sort(sort).
				All(&customers)
		}); err != nil {
		return customers, total, fmt.Errorf("Error[%s] while getting all customers", err)
	}
	return customers, total, nil
}

// Delete : delete customer by its ID
func (customer *Customer) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", CustomerCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": customer.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting customer by id[%s]", err, customer.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (customer *Customer) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": customer.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"name":     customer.Name,
			"email":    customer.Email,
			"phone":    customer.Phone,
			"gender":   customer.Gender,
			"birthday": customer.Birthday,
			"status":   customer.Status,

			"u_by": customer.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	return customer.Update(tenantID, selector, updator)
}
