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
	dbh "riverlife/internal/common/dbhandler"
	cmtypes "riverlife/internal/common/types"
	rlctypes "riverlife/internal/rlcollector/types"
	"sync"
	"time"
)

func ScheduleStateCollection(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(72 * time.Hour) // Run ever 3 days
	childCtx, cancel := context.WithCancel(ctx)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			rlctypes.Ctx.Log.Info("Shutting down State Collector")
			ticker.Stop()
			cancel()
			return
		case <-ticker.C:
			rlctypes.Ctx.Log.Info("Starting collection of state RSS feeds")
			doStateCollection(childCtx)
		}
	}
}

func doStateCollection(ctx context.Context) {
	stateJobs := make(chan rlctypes.State, 60)
	siteJobs := make(chan cmtypes.Site, 10000)
	var statewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var stats rlctypes.SafeCount

	for stw := 1; stw <= 10; stw++ {
		statewg.Add(1)
		go statesWorker(ctx, stw, stateJobs, siteJobs, &statewg, &stats)
	}

	for pw := 1; pw <= 10; pw++ {
		persistwg.Add(1)
		go persistentSiteWorker(ctx, pw, siteJobs, &persistwg, rlctypes.Ctx.DB, true)
	}

	states, err := loadStates()
	if err != nil {
		rlctypes.Ctx.Log.Fatal(err)
	}

	for _, state := range states {
		stateJobs <- state
	}
	close(stateJobs)
	statewg.Wait()

	close(siteJobs)
	persistwg.Wait()

	rlctypes.Ctx.Log.Printf("Found %d locations", stats.Count)
}

func ScheduleSiteCollection(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(1 * time.Hour) // Run ever 1 hours
	childCtx, cancel := context.WithCancel(ctx)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			rlctypes.Ctx.Log.Info("Shutting down Site Collector")
			ticker.Stop()
			cancel()
			return
		case <-ticker.C:
			rlctypes.Ctx.Log.Info("Starting collection of site xml")
			doSiteCollection(childCtx)
		}
	}
}

func doSiteCollection(ctx context.Context) {
	sites, err := dbh.GetSites(rlctypes.Ctx.DB)
	if err != nil {
		rlctypes.Ctx.Log.Error("Error retrieving sites from DB")
		rlctypes.Ctx.Log.Error(err)
		return
	}

	siteJobs := make(chan cmtypes.Site, 10000)
	persistJobs := make(chan cmtypes.Site, 10000)
	var sitewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var stats rlctypes.SafeCount

	for siw := 1; siw <= 10; siw++ {
		sitewg.Add(1)
		go siteWorker(ctx, siw, siteJobs, persistJobs, &sitewg, &stats)
	}

	for pw := 1; pw <= 10; pw++ {
		persistwg.Add(1)
		go persistentSiteWorker(ctx, pw, persistJobs, &persistwg, rlctypes.Ctx.DB, false)
	}

	for _, site := range sites {
		siteJobs <- *site
	}
	close(siteJobs)
	sitewg.Wait()

	close(persistJobs)
	persistwg.Wait()

	rlctypes.Ctx.Log.Printf("Found %d locations", stats.Count)
}

func DoStartupCollection(ctx context.Context, wg *sync.WaitGroup) {
	stateJobs := make(chan rlctypes.State, 60)
	siteJobs := make(chan cmtypes.Site, 10000)
	persistJobs := make(chan cmtypes.Site, 10000)
	var statewg sync.WaitGroup
	var sitewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var siteStats rlctypes.SafeCount
	var dataStats rlctypes.SafeCount

	defer wg.Done()

	for stw := 1; stw <= 10; stw++ {
		statewg.Add(1)
		go statesWorker(ctx, stw, stateJobs, siteJobs, &statewg, &siteStats)
	}

	for siw := 1; siw <= 10; siw++ {
		sitewg.Add(1)
		go siteWorker(ctx, siw, siteJobs, persistJobs, &sitewg, &dataStats)
	}

	for pw := 1; pw <= 10; pw++ {
		persistwg.Add(1)
		go persistentSiteWorker(ctx, pw, persistJobs, &persistwg, rlctypes.Ctx.DB, true)
	}

	states, err := loadStates()
	if err != nil {
		rlctypes.Ctx.Log.Fatal(err)
	}

	for _, state := range states {
		stateJobs <- state
	}

	close(stateJobs)
	rlctypes.Ctx.Log.Infof("Closed state job channel")
	statewg.Wait()

	close(siteJobs)
	rlctypes.Ctx.Log.Infof("Closed site job channel")
	sitewg.Wait()

	close(persistJobs)
	rlctypes.Ctx.Log.Infof("Closed persist job channel")
	persistwg.Wait()

	rlctypes.Ctx.Log.Infof("Found %d locations", siteStats.Count)
	rlctypes.Ctx.Log.Infof("Downloaded %f MB", float64(dataStats.Count/(1024*1024)))
}
