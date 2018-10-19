package model

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongo "nail/backend/mongojuice"
)

// Staff struct
type Staff struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Gender   string `json:"gender" bson:"gender"`
	Birthday string `json:"birthday" bson:"birthday"`

	AmountOwn  float64 `json:"amount_own" bson:"amount_own"`
	Salary     float64 `json:"salary" bson:"salary"`
	Commission float64 `json:"commission" bson:"commission"`

	Status int `json:"status" bson:"status"`

	CreatedBy string    `json:"c_by" bson:"c_by"`
	UpdatedBy string    `json:"u_by" bson:"u_by"`
	CreatedAt time.Time `json:"c_at" bson:"c_at"`
	UpdatedAt time.Time `json:"u_at" bson:"u_at"`
}

// StaffCollectionName returns the full name of staff collection
func StaffCollectionName(tenantID string) string {
	return fmt.Sprintf("%s_staff", tenantID)
}

// Insert : insert new staff record
func (staff *Staff) Insert(tenantID string) error {
	staff.ID = bson.NewObjectId().Hex()
	staff.UpdatedBy = staff.CreatedBy
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", StaffCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(staff)
		}); err != nil {
		return fmt.Errorf("Error[%s] while inserting staff", err)
	}

	return nil
}

// Update : udpate staff based on selector and updator
func (staff *Staff) Update(tenantID string, selector, updator interface{}) error {
	if err := mongo.Execute("monotonic", StaffCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Update(selector, updator)
		}); err != nil {
		return err
	}
	return nil
}

// FindByID : Find staff by its ID
func (staff *Staff) FindByID(tenantID string) error {
	if err := mongo.Execute("monotonic", StaffCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			selector := bson.M{
				"_id": staff.ID,
			}
			return collection.Find(selector).One(staff)
		}); err != nil {
		return fmt.Errorf("Error[%s] while getting staff by id[%s]", err, staff.ID)
	}

	return nil
}

// AllStaffs : find all staffs
func AllStaffs(tenantID string, filterInfo FilterInfo) ([]Staff, int, error) {
	total := 0
	var staffs []Staff
	if err := mongo.Execute("monotonic", StaffCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			total, _ = collection.Find(nil).Count()
			if (filterInfo == FilterInfo{}) {
				return collection.Find(nil).All(&staffs)
			}
			sort := filterInfo.PrepareSort()
			return collection.Find(nil).
				Limit(filterInfo.PerPage).
				Skip((filterInfo.Page - 1) * filterInfo.PerPage).
				Sort(sort).
				All(&staffs)
		}); err != nil {
		return staffs, total, fmt.Errorf("Error[%s] while getting all staffs", err)
	}
	return staffs, total, nil
}

// Delete : delete staff by its ID
func (staff *Staff) Delete(tenantID string) error {
	if err := mongo.Execute("monotonic", StaffCollectionName(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"_id": staff.ID,
			})
		}); err != nil {
		return fmt.Errorf("Error[%s] while deleting staff by id[%s]", err, staff.ID)
	}

	return nil
}

// UpdateAll : update all properties
func (staff *Staff) UpdateAll(tenantID string) error {
	selector := bson.M{
		"_id": staff.ID,
	}
	updator := bson.M{
		"$set": bson.M{
			"name":       staff.Name,
			"email":      staff.Email,
			"phone":      staff.Phone,
			"gender":     staff.Gender,
			"birthday":   staff.Birthday,
			"amount_own": staff.AmountOwn,
			"salary":     staff.Salary,
			"commission": staff.Commission,
			"u_by":       staff.UpdatedBy,
			"u_at":       time.Now(),
		},
	}
	return staff.Update(tenantID, selector, updator)
}
