package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

// Service struct
type Service struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// ServiceCollectionName returns the full name of service collection
func ServiceCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_service", tenantID)
}

// Insert : insert new service record
func (service *Service) Insert(tenantID string) error {
	service.ID = bson.NewObjectId().Hex()
	service.UpdatedBy = service.CreatedBy
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", ServiceCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(service)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting service", err)
	}

	return nil
}

// Update : udpate service based on selector and updator
func (service *Service) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", ServiceCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find service by its ID
func (service *Service) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", ServiceCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": service.ID,
			}
			return collection.Find(selector).One(service)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting service by id[%s]", err, service.ID)
	}

	return nil
}

// AllServices : find all services
func AllServices(tenantID string, filterInfo FilterInfo) ([]*Service, int, error) {
	total := 0
	var services []*Service
	if err := mongo.Execute("monotonic", ServiceCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			total, _ = collection.Find(nil).Count()
			if (filterInfo == FilterInfo{}) {
				return collection.Find(nil).All(&services)
			}
			sort := filterInfo.PrepareSort()
			return collection.Find(nil).
				Limit(filterInfo.PerPage).
				Skip((filterInfo.Page - 1) * filterInfo.PerPage).
				Sort(sort).
				All(&services)
		}); err != nil {
		return services, total, fmt.Errorf("Error[%s] while getting all services", err)
	}
	return services, total, nil
}

// Delete : delete service by its ID
func (service *Service) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", ServiceCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": service.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting service by id[%s]", err, service.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (service *Service) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": service.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"name": service.Name,

			"u_by": service.UpdatedBy,
			"u_at": time.Now(),
		},
	}
	return service.Update(tenantID, selector, updator)
}
