package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

// Tenant struct
type Tenant struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	Phone       string `json:"phone" bson:"phone"`
	AddressID   string `json:"address_id" bson:"address_id"`
	CompanyName string `json:"company_name" bson:"company_name"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// TenantCollectionName returns the full name of tenant collection
func TenantCollectionName() string {
	return fmt.Sprintf("tenant")
}

// Insert : insert new tenant record
func (tenant *Tenant) Insert() error {
	tenant.ID = bson.NewObjectId().Hex()
	tenant.UpdatedBy = tenant.CreatedBy
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", TenantCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Insert(tenant)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting tenant", err)
	}

	return nil
}

// Update : udpate tenant based on selector and updator
func (tenant *Tenant) Update(selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", TenantCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find tenant by its ID
func (tenant *Tenant) FindByID() error {
	if err := mongo.Execute("monotonic", TenantCollectionName(),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": tenant.ID,
			}
			return collection.Find(selector).One(tenant)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting tenant by id[%s]", err, tenant.ID)
	}

	return nil
}

// AllTenants : find all tenants
func AllTenants() ([]Tenant, error) {
	var tenants []Tenant
	if err := mongo.Execute("monotonic", TenantCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Find(nil).All(&tenants)
		}); err != nil {
		return tenants, fmt.Errorf("Error[%s] while getting all tenants", err)
	}
	return tenants, nil
}

// Delete : delete tenant by its ID
func (tenant *Tenant) Delete() error {
	if err := mongo.Execute("monotonic", TenantCollectionName(),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": tenant.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting tenant by id[%s]", err, tenant.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (tenant *Tenant) UpdateAll() error {
	selector := bson.M{
		"_id": tenant.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"name":         tenant.Name,
			"email":        tenant.Email,
			"phone":        tenant.Phone,
			"address_id":   tenant.AddressID,
			"company_name": tenant.CompanyName,
			"u_by":         tenant.UpdatedBy,
			"u_at":         time.Now(),
		},
	}
	return tenant.Update(selector, updator)
}
