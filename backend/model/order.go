package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

const (
	Booked     = 0
	InProgress = 1
	Finished   = 2
)

// Package struct
type Package struct {
	StaffID     string  `json:"staff_id" bson:"staff_id"`
	StaffName   string  `json:"staff_name" bson:"staff_name"`
	ServiceID   string  `json:"service_id" bson:"service_id"`
	ServiceName string  `json:"service_name" bson:"service_name"`
	ProductID   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Price       float64 `json:"price" bson:"price"`
	Duration    int     `json:"duration" bson:"duration"` // in minute
}

// Order struct
type Order struct {
	ID            string    `json:"id" bson:"_id"`
	CustomerID    string    `json:"customer_id" bson:"customer_id"`
	Customer      Customer  `json:"customer" bson:"-"`
	Packages      []Package `json:"packages" bson:"packages"`
	TotalPrice    float64   `json:"total_price" bson:"total_price"`
	TotalDuration int       `json:"total_duration" bson:"total_duration"` // in minute
	Status        int       `json:"status" bson:"status"`
	StartedAt     time.Time `json:"started_at" bson:"started_at"`
	EndedAt       time.Time `json:"ended_at" bson:"ended_at"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// OrderCollectionName returns the full name of order collection
func OrderCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_order", tenantID)
}

// Insert : insert new order record
func (order *Order) Insert(tenantID string) error {
	order.ID = bson.NewObjectId().Hex()
	order.Status = InProgress
	order.StartedAt = time.Now()
	order.EndedAt = order.StartedAt.Local().Add(time.Minute * time.Duration(order.TotalDuration))

	order.UpdatedBy = order.CreatedBy
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", OrderCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(order)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting order", err)
	}

	return nil
}

// Update : udpate order based on selector and updator
func (order *Order) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", OrderCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find order by its ID
func (order *Order) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", OrderCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": order.ID,
			}
			return collection.Find(selector).One(order)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting order by id[%s]", err, order.ID)
	}

	return nil
}

// AllOrders : find all orders
func AllOrders(tenantID string, status int) ([]*Order, error) {
	var orders []*Order
	selector := bson.M{}
	if status != -1 {
		selector = bson.M{
			"status": status,
		}
	}
	if err := mongo.Execute("monotonic", OrderCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Find(selector).All(&orders)
		}); err != nil {
		return orders, fmt.Errorf("Error[%s] while getting all orders", err)
	}
	for _, order := range orders {
		order.Customer.ID = order.CustomerID
		order.Customer.FindByID(tenantID)
	}
	return orders, nil
}

// Delete : delete order by its ID
func (order *Order) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", OrderCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": order.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting order by id[%s]", err, order.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (order *Order) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": order.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"customer_id":    order.CustomerID,
			"packages":       order.Packages,
			"total_price":    order.TotalPrice,
			"total_duration": order.TotalDuration,
			"status":         order.Status,
			"started_at":     order.StartedAt,
			"ended_at":       order.EndedAt,

			"u_by": order.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	return order.Update(tenantID, selector, updator)
}

// Checkout : checkout order
func (order *Order) Checkout(tenantID string, payment Payment) error {
	// update Order
	selector := bson.M{
		"_id": order.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"status": Finished,

			"u_by": order.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	order.Update(tenantID, selector, updator)

	// create billing
	billing := Billing{
		CustomerID:    order.CustomerID,
		OrderID:       order.ID,
		InitialAmount: order.TotalPrice,
		Surcharge:     payment.Surcharge,
		Discount:      payment.Discount,
		PaymentType:   payment.PaymentType,
		Status:        Paid,
		CreatedBy:     order.UpdatedBy,
	}
	// calculate TIP
	if payment.TipType == Share {
		billing.Tip = payment.Tip
		for _, staffTip := range payment.StaffTips {
			for _, pack := range order.Packages {
				if staffTip.StaffID == pack.StaffID {
					staffTip.MakeAmount += pack.Price
				}
			}
			staffTip.MakePercent = staffTip.MakeAmount / order.TotalPrice
			staffTip.TipAmount = staffTip.MakePercent * payment.Tip
			staff := &Staff{ID: staffTip.StaffID}
			if err := staff.FindByID(tenantID); err != nil {
				return fmt.Errorf("Error[%s] while calculating staff tip with order id[%s]", err, order.ID)
			}
			staff.AmountOwn += staffTip.TipAmount
			selector := bson.M{
				"_id": staff.ID,
			}
			updator := bson.M{
				"$set": bson.M{
					"amount_own": staff.AmountOwn,
				},
			}
			staff.Update(tenantID, selector, updator)
		}
	} else {
		for _, staffTip := range payment.StaffTips {
			billing.Tip += staffTip.TipAmount
			staff := &Staff{ID: staffTip.StaffID}
			if err := staff.FindByID(tenantID); err != nil {
				return fmt.Errorf("Error[%s] while calculating staff tip with order id[%s]", err, order.ID)
			}
			staff.AmountOwn += staffTip.TipAmount
			selector := bson.M{
				"_id": staff.ID,
			}
			updator := bson.M{
				"$set": bson.M{
					"amount_own": staff.AmountOwn,
				},
			}
			staff.Update(tenantID, selector, updator)
		}
	}
	billing.Revenue = billing.InitialAmount + billing.Surcharge - billing.Discount
	billing.TotalAmount = billing.Revenue + billing.Tip

	if err := billing.Insert(tenantID); err != nil {
		return fmt.Errorf("Error[%s] while checkout order by id[%s]", err, order.ID)
	}

	return nil
}
