package models

import (
	"time"

	"github.com/emanuelhristea/lime/config"
	"github.com/jinzhu/gorm"
)

// Subscription is a ...
type Subscription struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	StripeID   string    `gorm:"size:18;not null;unique" json:"stripe_id"`
	CustomerID uint64    `sql:"unique_index:idx_customer_tariff;type:int REFERENCES customers(id) ON DELETE CASCADE" json:"customer_id"`
	TariffID   uint64    `sql:"unique_index:idx_customer_tariff;type:int REFERENCES tariffs(id) ON DELETE CASCADE" json:"tariff_id"`
	Status     bool      `gorm:"false" json:"status"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Customer   Customer  `json:"customer"`
	Tariff     Tariff    `json:"tariff"`
	IssuedAt   time.Time `json:"issued_at"`  // Issued At
	ExpiresAt  time.Time `json:"expires_at"` // Expires At
	Licenses   []License `json:"licenses"`
}

// Expired is a ...
func (s *Subscription) Expired() bool {
	return !s.IssuedAt.IsZero() && time.Now().After(s.ExpiresAt)
}

// SaveSubscription is a ...
func (s *Subscription) SaveSubscription() (*Subscription, error) {
	err := config.DB.Create(&s).Error
	if err != nil {
		return &Subscription{}, err
	}
	return s, nil
}

// FindSubscriptionByID is a ...
func (s *Subscription) FindSubscriptionByID(uid uint64, relations ...string) (*Subscription, error) {
	db := config.DB.Model(Subscription{}).Where("id = ?", uid)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}

	err := db.Take(&s).Error
	if err != nil {
		return &Subscription{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Subscription{}, ErrTariffNotFound
	}
	return s, err
}

// FindSubscriptionByStripeID is a ...
func (s *Subscription) FindSubscriptionByStripeID(stripeID string, relations ...string) (*Subscription, error) {
	db := config.DB.Model(Subscription{}).Where("stripe_id = ? AND status = ?", stripeID, true)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}

	err := db.Take(&s).Error
	if err != nil {
		return &Subscription{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Subscription{}, ErrCustomerNotFound
	}
	return s, err
}

// UpdateSubscription is a ...
func (s *Subscription) UpdateSubscription(uid uint64) (*Subscription, error) {
	db := config.DB.Model(&Subscription{}).Where("id = ?", uid).Take(&Subscription{}).UpdateColumns(
		map[string]interface{}{
			"customer_id": s.CustomerID,
			"stripe_id":   s.StripeID,
			"tariff_id":   s.TariffID,
			"status":      s.Status,
			"issued_at":   s.IssuedAt,
			"expires_at":  s.ExpiresAt,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return &Subscription{}, db.Error
	}

	err := db.Model(&Subscription{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Subscription{}, err
	}
	return s, nil
}

// DeleteSubscription is a ...
func DeleteSubscription(uid uint64) (int64, error) {
	db := config.DB.Model(&Subscription{}).Where("id = ?", uid).Take(&Subscription{}).Delete(&Subscription{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// CustomerSubscriptionsList is a ...
type CustomerSubscriptionsList struct {
	ID           uint64    `json:"id"`
	StripeID     string    `json:"stripe_id"`
	CustomerID   string    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	TariffID     uint64    `json:"tariff_id"`
	TariffName   string    `json:"tariff_name"`
	Status       bool      `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SubscriptionsByCustomerID is a ...
func SubscriptionsByCustomerID(customerID string) *[]CustomerSubscriptionsList {
	result := []CustomerSubscriptionsList{}
	db := config.DB.Raw("SELECT subscriptions.ID,subscriptions.stripe_id,subscriptions.customer_id,customers.NAME AS customer_name,subscriptions.tariff_id,tariffs.NAME AS tariff_name,subscriptions.status,subscriptions.created_at,subscriptions.updated_at FROM subscriptions INNER JOIN customers ON subscriptions.customer_id=customers.ID INNER JOIN tariffs ON subscriptions.tariff_id=tariffs.ID WHERE subscriptions.customer_id=? ORDER BY subscriptions.created_at DESC ", customerID).Scan(&result)
	if db.Error != nil {
		return &result
	}
	return &result
}

// SubscriptionsList is a ...
func SubscriptionsList(customerID string, relations ...string) *[]Subscription {
	db := config.DB.Model(&Subscription{}).Where("customer_id=?", customerID).Order("ID asc")
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}

	subscriptions := []Subscription{}
	db = db.Find(&subscriptions)
	if db.Error != nil {
		return &subscriptions
	}
	return &subscriptions
}
