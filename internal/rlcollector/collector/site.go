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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	cmtypes "riverlife/internal/common/types"
	rlctypes "riverlife/internal/rlcollector/types"
	"strconv"
	"sync"
	"time"
)

func siteWorker(ctx context.Context, id int,
	siteJobs <-chan cmtypes.Site,
	persistJobs chan<- cmtypes.Site,
	wg *sync.WaitGroup,
	sc *rlctypes.SafeCount) {
	rlctypes.Ctx.Log.Infof("Starting Site Worker %d", id)
	defer wg.Done()
	for site := range siteJobs {
		select {
		case <-ctx.Done():
			rlctypes.Ctx.Log.Infof("Site Worker %d interrupted", id)
			rlctypes.Ctx.Log.Infof("Shutting down Site Worker %d", id)
			return
		default:
			rlctypes.Ctx.Log.Infof("Starting parser on Site Worker %d", id)
			sc.IncCount(parseSiteXML(site, persistJobs))
			rlctypes.Ctx.Log.Infof("Ending parser on Site Worker %d", id)
		}
	}
	rlctypes.Ctx.Log.Infof("Site Worker %d finished", id)
}

func parseSiteXML(site cmtypes.Site, persistJobs chan<- cmtypes.Site) int64 {
	// Skip stations that are out of service or not collecting data
	if site.IsInService == false || site.HasData == false {
		rlctypes.Ctx.Log.Infof("Site %s (%s) has no data at this time, skipping data collection.", site.ID, site.Location)
		return 0
	}

	rlctypes.Ctx.Log.Infof("Loading site data for " + site.ID)
	req, err := http.NewRequest("GET", fmt.Sprintf(rlctypes.SiteURL, site.ID), nil)
	if err != nil {
		rlctypes.Ctx.Log.Errorf("Error finding site " + site.ID)
		rlctypes.Ctx.Log.Error(err)
		return 0
	}
	bytes, err := doRequest(req)
	if err != nil {
		rlctypes.Ctx.Log.Errorf("Error finding site " + site.ID)
		rlctypes.Ctx.Log.Error(err)
		return 0
	}

	var noaaSite rlctypes.NoaaSite
	err = xml.Unmarshal(bytes, &noaaSite)
	if err != nil {
		rlctypes.Ctx.Log.Errorf("Error finding site " + site.ID)
		rlctypes.Ctx.Log.Error(err)
		return 0
	}
	if len(noaaSite.Observed.Datums) > 0 {
		site.CurrentLevel, err = strconv.ParseFloat(noaaSite.Observed.Datums[0].Primary.Text, 64)
		if err != nil {
			rlctypes.Ctx.Log.Warnf("Error parsing staging value for site %s", site.ID)
			rlctypes.Ctx.Log.Warn(err)
			site.CurrentLevel = 0.0
		}

		if noaaSite.Observed.Datums[0].Secondary.Text != "" {
			site.CurrentFlow, err = strconv.ParseFloat(noaaSite.Observed.Datums[0].Secondary.Text, 64)
			if err != nil {
				rlctypes.Ctx.Log.Warnf("Error parsing flow value for site %s", site.ID)
				rlctypes.Ctx.Log.Warn(err)
				site.CurrentFlow = 0.0
			}
		} else {
			site.CurrentFlow = 0.0
		}

		if site.CurrentFlow < 0 {
			site.CurrentFlow = 0.0
		}

		site.SampleTime, err = time.Parse(time.RFC3339, noaaSite.Observed.Datums[0].Valid.Text)
		if err != nil {
			rlctypes.Ctx.Log.Warnf("Error parsing time value for site %s", site.ID)
			rlctypes.Ctx.Log.Warn(err)
		}
	}
	persistJobs <- site
	rlctypes.Ctx.Log.Infof("Finished loading site for " + site.ID)
	return int64(len(bytes))
}

func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}
