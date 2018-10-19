package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

// Product struct
type Product struct {
	ID        string  `json:"id" bson:"_id"`
	Name      string  `json:"name" bson:"name"`
	ServiceID string  `json:"service_id" bson:"service_id"`
	Price     float64 `json:"price" bson:"price"`
	Duration  int     `json:"duration" bson:"duration"` // in minute

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"-" bson:"c_at"`
	UpdatedAt time.Time `json:"-" bson:"u_at"`
}

// ProductCollectionName returns the full name of product collection
func ProductCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_product", tenantID)
}

// Insert : insert new product record
func (product *Product) Insert(tenantID string) error {
	product.ID = bson.NewObjectId().Hex()
	product.UpdatedBy = product.CreatedBy
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(product)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting product", err)
	}

	return nil
}

// Update : udpate product based on selector and updator
func (product *Product) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find product by its ID
func (product *Product) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": product.ID,
			}
			return collection.Find(selector).One(product)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting product by id[%s]", err, product.ID)
	}

	return nil
}

// FindProductsByServiceID : find products by service ID
func FindProductsByServiceID(tenantID, serviceID string, filterInfo FilterInfo) ([]*Product, int, error) {
	total := 0
	var products []*Product
	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"service_id": serviceID,
			}
			total, _ = collection.Find(selector).Count()
			if (filterInfo == FilterInfo{}) {
				return collection.Find(selector).All(&products)
			}
			sort := filterInfo.PrepareSort()
			return collection.Find(nil).
				Limit(filterInfo.PerPage).
				Skip((filterInfo.Page - 1) * filterInfo.PerPage).
				Sort(sort).
				All(&products)
		}); err != nil {
		return products, total, fmt.Errorf("Error[%s] while finding products by service_id[%s]", err, serviceID)
	}
	return products, total, nil
}

// AllProducts : find all products
func AllProducts(tenantID string) ([]*Product, error) {
	var products []*Product
	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Find(nil).All(&products)
		}); err != nil {
		return products, fmt.Errorf("Error[%s] while getting all products", err)
	}
	return products, nil
}

// Delete : delete product by its ID
func (product *Product) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", ProductCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": product.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting product by id[%s]", err, product.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (product *Product) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": product.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"name":       product.Name,
			"service_id": product.ServiceID,
			"price":      product.Price,
			"duration":   product.Duration,

			"u_by": product.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	return product.Update(tenantID, selector, updator)
}
