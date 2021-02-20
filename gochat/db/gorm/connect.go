package gorm

import "github.com/jinzhu/gorm"

// Connect to database with gorm.
func Connect(connInfo string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connInfo)

	if err != nil {
		return nil, err
	}
	return db, nil
}
