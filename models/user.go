package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string `gorm:"type:varchar(255);not null;unique" json:"name"`
	Password        string `gorm:"type:varchar(255);not null" json:"password"`
	FollowCount     int64  `gorm:"type:int;not null" json:"follow_count"`
	FollowerCount   int64  `gorm:"type:int;not null" json:"follower_count"`
	IsFollow        bool   `gorm:"type:bool;not null" json:"is_follow"`
	Avatar          string `gorm:"type:varchar(255);not null" json:"avatar"`
	BackgroundImage string `gorm:"type:varchar(255);not null" json:"background_image"`
	Signature       string `gorm:"type:varchar(255);not null" json:"signature"`
	TotalFavorited  int64  `gorm:"type:int;not null" json:"total_favorited"`
	WorkCount       int64  `gorm:"type:int;not null" json:"work_count"`
	FavoriteCount   int64  `gorm:"type:int;not null" json:"favorite_count"`
}

// BeforeCreate is a hook that will be called before creating a new user.
// It hashes the password using bcrypt.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Password = string(hashedPassword)

	return nil
}

// CreateNewUser creates a new user.
// It returns the user if it is created successfully, otherwise it returns an error.
//
//	@param u *User
//	@return *User
//	@return error
func CreateNewUser(u *User) (*User, error) {
	err := db.Model(&User{}).Create(u).Error

	return u, err
}

// GetUserByName gets a user by its name.
// It returns the user if it is found, otherwise it returns an error.
//
//	@param name string
//	@return *User
//	@return error
func GetUserByName(name string) (*User, error) {
	user := &User{}
	err := db.Model(&User{}).Where("name = ?", name).First(user).Error

	return user, err
}

// GetUserByUserID gets a user by its id.
//
//	@param id int64
//	@return *User
func GetUserByUserID(id int64) (*User, error) {
	user := &User{}
	err := db.Model(&User{}).Where("id = ?", id).First(user).Error

	return user, err
}
