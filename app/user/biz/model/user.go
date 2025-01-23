package model

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(20);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}

func (User) TableName() string {
	return "users"
}

func GetUserById(db *gorm.DB, ctx context.Context, id uint32) (*User, error) {
	user := &User{}
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Model: gorm.Model{ID: uint(id)}}).First(user).Error
	return user, err
}

func GetUserByUsername(db *gorm.DB, ctx context.Context, username string) (*User, error) {
	user := &User{}
	err := db.WithContext(ctx).Model(&User{}).Where(&User{Username: username}).First(user).Error
	return user, err
}

func CreateUser(db *gorm.DB, ctx context.Context, username, passwordHash string) (*User, error) {
	user := &User{Username: username, PasswordHash: passwordHash}
	return user, db.WithContext(ctx).Model(&User{}).Create(user).Error
}

func DeleteUser(db *gorm.DB, ctx context.Context, id uint) error {
	return db.WithContext(ctx).Model(&User{}).Delete(&User{Model: gorm.Model{ID: id}}).Error
}

func UpdateUser(db *gorm.DB, ctx context.Context, id uint, passwordHash string) error {
	return db.WithContext(ctx).Model(&User{Model: gorm.Model{ID: id}}).Updates(&User{PasswordHash: passwordHash}).Error
}
