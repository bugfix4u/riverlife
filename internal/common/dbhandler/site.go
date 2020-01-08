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
package dbhandler

import (
	"github.com/jinzhu/gorm"
	cmtypes "riverlife/internal/common/types"
)

func GetSites(db *gorm.DB) ([]*cmtypes.Site, error) {
	var sites []*cmtypes.Site
	result := db.Find(&sites)
	err := DbResults(result)

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return sites, nil
}

func CreateorUpdateSite(db *gorm.DB, site *cmtypes.Site, checkCreate bool) error {
	var result *gorm.DB
	if checkCreate {
		result = db.Where("id = ?", site.ID).Assign(site).FirstOrCreate(site)
	} else {
		result = db.Update(site)
	}
	return DbResults(result)
}

func GetSite(db *gorm.DB, site *cmtypes.Site) error {
	result := db.Where("id = ?", site.ID).Find(site)
	return DbResults(result)
}
