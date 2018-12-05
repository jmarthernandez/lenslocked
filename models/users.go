package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// needed for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

// ById will look up User by id provided
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will look up User by email provided
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err

}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}
}

// Create will create a User and return the ID
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) Create(user *User) error {
	return us.db.Create(&user).Error
}

// Delete will delete a User
func (us *UserService) Delete(id uint) error {
	// stupid gorm will drop the entire table if id is 0
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close will close UserService db connection
// backfill data id, createdAt, updatedAt
func (us *UserService) Close() error {
	return us.db.Close()
}

// Update will updates the user with all provided data
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// DesctructiveReset drops the user table and rebuilds it
func (us *UserService) DesctructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automigrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
