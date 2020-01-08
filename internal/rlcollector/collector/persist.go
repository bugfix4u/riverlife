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
package collector

import (
	"context"
	"sync"
	"github.com/jinzhu/gorm"
	cmtypes "riverlife/internal/common/types"
	dbh "riverlife/internal/common/dbhandler"
	rlctypes "riverlife/internal/rlcollector/types"
)

func persistentSiteWorker(ctx context.Context, id int, siteJobs <-chan cmtypes.Site, wg *sync.WaitGroup, db *gorm.DB, checkCreate bool) {
	rlctypes.Ctx.Log.Infof("Starting Persistent Worker %d", id)
	defer wg.Done()
	for site := range siteJobs {
		select {
		case <-ctx.Done():
			rlctypes.Ctx.Log.Infof("Persistent Worker %d interrupted", id)
			rlctypes.Ctx.Log.Infof("Shutting down Persistent Worker %d", id)
			return
		default:
			rlctypes.Ctx.Log.Infof("Updating database for site %s", site.ID)
			if err := dbh.CreateorUpdateSite(db, &site, checkCreate); err != nil {
				rlctypes.Ctx.Log.Errorf("Failed to update DB with site %s", site.ID)
			}
		}
	}
	rlctypes.Ctx.Log.Infof("Persistent Worker %d finished", id)
}