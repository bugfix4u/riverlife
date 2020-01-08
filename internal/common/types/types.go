// River Life
// Copyright (C) 2020  Denny Chsmbers

// This progrsm is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This progrsm is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this progrsm.  If not, see <http://www.gnu.org/licenses/>.
package types

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ActionType string

const (
	ActionTypeUnknown  ActionType = "Unknown"
	ActionTypeNone     ActionType = "None"
	ActionTypeAction   ActionType = "Action"
	ActionTypeMinor    ActionType = "Minor"
	ActionTypeModerate ActionType = "Moderate"
	ActionTypeMajor    ActionType = "Major"
	ActionTypeLowWater ActionType = "Low Water"
)

type Site struct {
	ID            string     `json:"id" gorm:"column:id;primary_key"`
	Location      string     `json:"location" gorm:"not null"`
	State         string     `json:"state"`
	IsCurrent     bool       `json:"isCurrent"`
	IsInService   bool       `json:"isInService"`
	HasData       bool       `json:"hasData"`
	CurrentLevel  float64    `json:"currentLevel"`
	CurrentFlow   float64    `json:"currentFlow"`
	SampleTime    time.Time  `json:"sampleTime"`
	CurrentAction ActionType `json:"currentAction"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

//GORM methods

// TableName alters the given name
func (site *Site) TableName() string {
	return "site"
}

// BeforeCreate is a GORM callback function to do operation before the create is called.
func (site *Site) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreatedAt", time.Now())
	if err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM callback function to do operation before the update is called.
func (site *Site) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdatedAt", time.Now())
	if err != nil {
		return err
	}

	return nil
}
