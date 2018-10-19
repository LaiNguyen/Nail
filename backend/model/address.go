package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

// Address struct
type Address struct {
	ID         string `json:"id" bson:"_id"`
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	Province   string `json:"province" bson:"province"`
	Country    string `json:"country" bson:"country"`
	PostalCode string `json:"postal_code" bson:"postal_code"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// AddressCollectionName returns the full name of address collection
func AddressCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_address", tenantID)
}

// Insert : insert new address record
func (address *Address) Insert(tenantID string) error {
	address.ID = bson.NewObjectId().Hex()
	address.UpdatedBy = address.CreatedBy
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", AddressCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(address)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting address", err)
	}

	return nil
}

// Update : udpate address based on selector and updator
func (address *Address) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", AddressCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find address by its ID
func (address *Address) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", AddressCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": address.ID,
			}
			return collection.Find(selector).One(address)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting address by id[%s]", err, address.ID)
	}

	return nil
}

// AllAddress : find all address
func AllAddresses(tenantID string) ([]Address, error) {
	var Addresses []Address
	if err := mongo.Execute("monotonic", AddressCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Find(nil).All(&Addresses)
		}); err != nil {
		return Addresses, fmt.Errorf("Error[%s] while getting all Addresses", err)
	}
	return Addresses, nil
}

// Delete : delete address by its ID
func (address *Address) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", AddressCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": address.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting address by id[%s]", err, address.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (address *Address) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": address.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"street":      address.Street,
			"city":        address.City,
			"province":    address.Province,
			"country":     address.Country,
			"postal_code": address.PostalCode,
			"u_by":        address.UpdatedBy,
			"u_at":        time.Now(),
		},
	}
	return address.Update(tenantID, selector, updator)
}
