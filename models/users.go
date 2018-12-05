package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// needed for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound = errors.New("models: resource not found")
)

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
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create a User and return the ID
// 1 - user, nil
// 2 - nil, ErrNotFound
// 3 - nil, otherError
func (us *UserService) Create(user *User) error {
	return us.db.Create(&user).Error
}

// Close closes UserService db connection
// backfill data id, createdAt, updatedAt
func (us *UserService) Close() error {
	return us.db.Close()
}

// DesctructiveReset drops the user table and rebuilds it
func (us *UserService) DesctructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unqiue_index"`
}
