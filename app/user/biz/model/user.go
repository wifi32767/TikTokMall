package model

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(20);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Permission   uint32 `gorm:"not null"`
	// 权限字段暂且分为两类
	// 1: 普通用户
	// 2: 管理员
	// 由于0是空值，不好区分，从1开始
	// 普通用户只能访问和操作自己所拥有的资源，可以访问商品
	// 管理员可以访问和操作所有资源，可以为其他账号赋予权限
	// 后续可能会扩展为更复杂的权限系统，比如通过k位四进制数管理k种不同的服务
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
	user := &User{Username: username, PasswordHash: passwordHash, Permission: 0}
	return user, db.WithContext(ctx).Model(&User{}).Create(user).Error
}

func DeleteUser(db *gorm.DB, ctx context.Context, id uint) error {
	return db.WithContext(ctx).Model(&User{}).Delete(&User{Model: gorm.Model{ID: id}}).Error
}

func UpdateUser(db *gorm.DB, ctx context.Context, id uint, passwordHash string) error {
	return db.WithContext(ctx).Model(&User{Model: gorm.Model{ID: id}}).Updates(&User{PasswordHash: passwordHash}).Error
}

func GrantUser(db *gorm.DB, ctx context.Context, id uint, permission uint32) error {
	return db.WithContext(ctx).Model(&User{Model: gorm.Model{ID: id}}).Updates(&User{Permission: permission}).Error
}
