// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"fmt"
	"time"

	"github.com/marmotedu/component-base/pkg/auth"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"gorm.io/gorm"
)

// User represents a user restful resource. It is also used as gorm model.
type User struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status int `json:"status" gorm:"column:status" validate:"omitempty"`

	// Required: true
	Nickname string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" gorm:"column:password" validate:"required"`

	// Required: true
	Email string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" gorm:"column:phone" validate:"omitempty"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`

	LoginedAt time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

// TableName maps to mysql table name.
func (u *User) TableName() string {
	return "user"
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}

	return nil
}

// AfterCreate run after create database record.
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}
