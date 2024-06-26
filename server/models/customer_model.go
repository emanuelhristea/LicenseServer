package models

import (
	"time"

	"github.com/emanuelhristea/lime/config"
	"github.com/emanuelhristea/lime/license"
	"github.com/jinzhu/gorm"
)

// Customer is a ...
type Customer struct {
	ID            uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Name          string         `gorm:"size:255;not null;unique" json:"name"`
	Email         string         `gorm:"size:255;not null;unique" json:"email"`
	Role          license.Role   `sql:"type:role" json:"role"`
	Status        bool           `gorm:"false" json:"status"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Subscriptions []Subscription `json:"subscriptions"`
}

func (c *Customer) IsAdmin() bool {
	return c.Role == license.AdminRole
}

func (c *Customer) IsUser() bool {
	return c.Role == license.UserRole
}

func (c *Customer) IsGuest() bool {
	return c.Role == license.GuestRole
}

// SaveCustomer is a ...
func (c *Customer) SaveCustomer() (*Customer, error) {
	err := config.DB.Create(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	return c, nil
}

// FindCustomerByID is a ...
func (c *Customer) FindCustomerByID(uid uint64, relations ...string) (*Customer, error) {
	db := config.DB.Model(Customer{}).Where("id = ?", uid)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}

	err := db.Take(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Customer{}, ErrCustomerNotFound
	}
	return c, err
}

// FindCustomerByEmail is a ...
func (c *Customer) FindCustomerByEmail(email string, relations ...string) (*Customer, error) {
	db := config.DB.Model(Customer{}).Where("email = ?", email).Where("status = ?", true)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}
	err := db.Take(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Customer{}, ErrCustomerNotFound
	}
	return c, err
}

// UpdateCustomer is a ...
func (c *Customer) UpdateCustomer(uid uint64) (*Customer, error) {
	db := config.DB.Model(&Customer{}).Where("id = ?", uid).Take(&Customer{}).UpdateColumns(
		map[string]interface{}{
			"name":       c.Name,
			"status":     c.Status,
			"email":      c.Email,
			"role":       c.Role,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Customer{}, db.Error
	}

	err := db.Model(&Customer{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Customer{}, err
	}
	return c, nil
}

// DeleteCustomer is a ...
func DeleteCustomer(uid uint64) (int64, error) {
	db := config.DB.Model(&Customer{}).Where("id = ?", uid).Take(&Customer{}).Delete(&Customer{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// CustomersList is a ...
func CustomersList(relations ...string) *[]Customer {
	db := config.DB.Model(&Customer{})
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}
	customers := []Customer{}
	db = db.Order("ID asc").Find(&customers)

	if db.Error != nil {
		return &customers
	}
	return &customers
}
