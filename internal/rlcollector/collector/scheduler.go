// River Life
// Copyright (C) 2020  Denny Chambers

// This program is free software: you can redistribute it and/or modify
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
	"fmt"
	redis "github.com/mediocregopher/radix/v3"
	"net/http"
	dbh "riverlife/internal/common/dbhandler"
	cmtypes "riverlife/internal/common/types"
	rlctypes "riverlife/internal/rlcollector/types"
	"sync"
	"time"
)

func ScheduleStateCollection(ctx context.Context, wg *sync.WaitGroup) {
	rlctypes.Ctx.Log.Infof("Setting up state collection timer for every %f hours", rlctypes.Ctx.Config.StateTickerTimeHour.Hours())
	ticker := time.NewTicker(rlctypes.Ctx.Config.StateTickerTimeHour)
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
	stateJobs := make(chan rlctypes.State, rlctypes.Ctx.Config.StateChannelSize)
	siteJobs := make(chan cmtypes.Site, rlctypes.Ctx.Config.SiteChannelSize)
	var statewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var siteStats rlctypes.SafeCount
	var dataStats rlctypes.SafeCount

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	for stw := 1; stw <= rlctypes.Ctx.Config.StateWorkerThreadCount; stw++ {
		statewg.Add(1)
		go statesWorker(ctx, stw, stateJobs, siteJobs, &statewg, httpClient, &siteStats, &dataStats)
	}

	for pw := 1; pw <= rlctypes.Ctx.Config.PersistWorkerThreadCount; pw++ {
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

	rlctypes.Ctx.Log.Infof("Found %d locations", siteStats.Count)
	rlctypes.Ctx.Log.Infof("Downloaded %f MB", float64(dataStats.Count)/float64(1024*1024))

}

func ScheduleSiteCollection(ctx context.Context, wg *sync.WaitGroup) {
	rlctypes.Ctx.Log.Infof("Setting up site collection timer for every %f hours", rlctypes.Ctx.Config.SiteTickerTimeHour.Hours())
	ticker := time.NewTicker(rlctypes.Ctx.Config.SiteTickerTimeHour)
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

	siteJobs := make(chan cmtypes.Site, rlctypes.Ctx.Config.SiteChannelSize)
	persistJobs := make(chan cmtypes.Site, rlctypes.Ctx.Config.PersistChannelSize)
	var sitewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var stats rlctypes.SafeCount

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	redisPool := createRedisClientPool()
	defer redisPool.Close()

	for siw := 1; siw <= rlctypes.Ctx.Config.SiteWorkerThreadCount; siw++ {
		sitewg.Add(1)
		go siteWorker(ctx, siw, siteJobs, persistJobs, &sitewg, redisPool, httpClient, &stats)
	}

	for pw := 1; pw <= rlctypes.Ctx.Config.PersistWorkerThreadCount; pw++ {
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

	rlctypes.Ctx.Log.Infof("Downloaded %f MB", float64(stats.Count)/float64(1024*1024))
}

func DoStartupCollection(ctx context.Context, wg *sync.WaitGroup) {
	stateJobs := make(chan rlctypes.State, rlctypes.Ctx.Config.StateChannelSize)
	siteJobs := make(chan cmtypes.Site, rlctypes.Ctx.Config.SiteChannelSize)
	persistJobs := make(chan cmtypes.Site, rlctypes.Ctx.Config.PersistChannelSize)
	var statewg sync.WaitGroup
	var sitewg sync.WaitGroup
	var persistwg sync.WaitGroup
	var siteStats rlctypes.SafeCount
	var dataStats rlctypes.SafeCount

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	redisPool := createRedisClientPool()
	defer redisPool.Close()

	defer wg.Done()

	for stw := 1; stw <= rlctypes.Ctx.Config.StateWorkerThreadCount; stw++ {
		statewg.Add(1)
		go statesWorker(ctx, stw, stateJobs, siteJobs, &statewg, httpClient, &siteStats, &dataStats)
	}

	for siw := 1; siw <= rlctypes.Ctx.Config.SiteWorkerThreadCount; siw++ {
		sitewg.Add(1)
		go siteWorker(ctx, siw, siteJobs, persistJobs, &sitewg, redisPool, httpClient, &dataStats)
	}

	for pw := 1; pw <= rlctypes.Ctx.Config.PersistWorkerThreadCount; pw++ {
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
	rlctypes.Ctx.Log.Infof("Downloaded %f MB", float64(dataStats.Count)/float64(1024*1024))
}

func createRedisClientPool() *redis.Pool {
	//Create a redis client pool with the same number of connections as the site worker thread count
	redisAddress := fmt.Sprintf("%s:%s", rlctypes.Ctx.Config.RedisHost, rlctypes.Ctx.Config.RedisPort)
	client, err := redis.NewPool("tcp", redisAddress, rlctypes.Ctx.Config.SiteWorkerThreadCount)
	if err != nil {
		rlctypes.Ctx.Log.Error("Error creating redis client pool")
		return nil
	}

	return client
}
