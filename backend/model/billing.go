package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

const (
	Unpaid = 0
	Paid   = 1
	Refund = 2
)

const (
	Share    = 0
	Specific = 1
)

// StaffTip struct
type StaffTip struct {
	StaffID     string  `json:"staff_id"`
	TipAmount   float64 `json:"tip_amount"`
	MakeAmount  float64 `json:"-"`
	MakePercent float64 `json:"-"`
}

// Payment struct
type Payment struct {
	Surcharge   float64    `json:"surcharge"`
	Discount    float64    `json:"discount"`
	PaymentType int        `json:"payment_type"`
	TipType     int        `json:"tip_type"`
	Tip         float64    `json:"tip"`
	StaffTips   []StaffTip `json:"staff_tips"`
}

// Billing struct
type Billing struct {
	ID            string  `json:"id" bson:"_id"`
	CustomerID    string  `json:"customer_id" bson:"customer_id"`
	OrderID       string  `json:"order_id" bson:"order_id"`
	InitialAmount float64 `json:"initial_amount" bson:"initial_amount"`
	Surcharge     float64 `json:"surcharge" bson:"surcharge"`
	Discount      float64 `json:"discount" bson:"discount"`
	Tip           float64 `json:"tip" bson:"tip"`
	TotalAmount   float64 `json:"total_amount" bson:"total_amount"`
	Revenue       float64 `json:"revenue" bson:"revenue"`
	Status        int     `json:"status" bson:"status"`
	PaymentType   int     `json:"payment_type" bson:"payment_type"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// BillingCollectionName returns the full name of billing collection
func BillingCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_billing", tenantID)
}

// Insert : insert new billing record
func (billing *Billing) Insert(tenantID string) error {
	billing.ID = bson.NewObjectId().Hex()

	billing.UpdatedBy = billing.CreatedBy
	billing.CreatedAt = time.Now()
	billing.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", BillingCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(billing)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting billing", err)
	}

	return nil
}

// Update : udpate billing based on selector and updator
func (billing *Billing) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", BillingCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find billing by its ID
func (billing *Billing) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", BillingCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": billing.ID,
			}
			return collection.Find(selector).One(billing)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting billing by id[%s]", err, billing.ID)
	}

	return nil
}

// AllBillings : find all billings
func AllBillings(tenantID string, filterInfo FilterInfo) ([]*Billing, int, error) {
	sort := filterInfo.PrepareSort()
	total := 0
	var billings []*Billing
	// selector := bson.M{}
	// if status != -1 {
	// 	selector = bson.M{
	// 		"status": status,
	// 	}
	// }
	if err := mongo.Execute("monotonic", BillingCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			// return collection.Find(nil).All(&billings)
			total, _ = collection.Find(nil).Count()
			return collection.Find(nil).
				Limit(filterInfo.PerPage).
				Skip((filterInfo.Page - 1) * filterInfo.PerPage).
				Sort(sort).
				All(&billings)
		}); err != nil {
		return billings, total, fmt.Errorf("Error[%s] while getting all billings", err)
	}
	return billings, total, nil
}

// Delete : delete billing by its ID
func (billing *Billing) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", BillingCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": billing.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting billing by id[%s]", err, billing.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (billing *Billing) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": billing.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"customer_id":    billing.CustomerID,
			"order_id":       billing.OrderID,
			"initial_amount": billing.InitialAmount,
			"surcharge":      billing.Surcharge,
			"discount":       billing.Discount,
			"tip":            billing.Tip,
			"total_amount":   billing.TotalAmount,
			"revenue":        billing.Revenue,
			"status":         billing.Status,

			"u_by": billing.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	return billing.Update(tenantID, selector, updator)
}
