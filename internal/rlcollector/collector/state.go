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
	"fmt"
	"strings"
	"sync"
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	rlctypes "riverlife/internal/rlcollector/types"
	cmtypes "riverlife/internal/common/types"
	"github.com/mmcdole/gofeed"
)

func statesWorker(ctx context.Context, id int, 
									stateJobs <-chan rlctypes.State,
									siteJobs chan<- cmtypes.Site,
									wg *sync.WaitGroup,
									count *rlctypes.SafeCount) {

	rlctypes.Ctx.Log.Infof("Starting State Worker %d", id)
	defer wg.Done()
	for state := range stateJobs {
		select {
		case <-ctx.Done():
			rlctypes.Ctx.Log.Infof("State Worker %d interrupted", id)
			rlctypes.Ctx.Log.Infof("Shutting down State Worker %d", id)
			return
		default:
			rlctypes.Ctx.Log.Infof("Starting parser on State Worker %d", id)
			count.IncCount(parseStateRSS(state, siteJobs))
			rlctypes.Ctx.Log.Infof("Ending parser on State Worker %d", id)
		}
	}		
	rlctypes.Ctx.Log.Infof("State Worker %d finished", id)
}

func parseStateRSS(state rlctypes.State, siteJobs chan<- cmtypes.Site) int64 {
rlctypes.Ctx.Log.Infof("Loading locations for " + state.State)
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(fmt.Sprintf(rlctypes.StateURL, state.Abbr))
		if err != nil {
			rlctypes.Ctx.Log.Errorf("Error finding state " + state.State)
			rlctypes.Ctx.Log.Error(err)
			return 0
		}
		var site cmtypes.Site
		var count int64
		for _, item := range feed.Items {
			count++

			stringTok := strings.Split(item.Title, "Observation -")
			var obsInfo string = ""
			if stringTok[0] != "" {
				loc := strings.LastIndex(stringTok[0], "-")
				if loc > -1 {
					runes := []rune(stringTok[0])
					obsInfo = string(runes[0:loc])
				} else {
					obsInfo = stringTok[0]
				}
				obsInfo = strings.TrimSpace(obsInfo)
				getObservationInfo(obsInfo, &site)
			} else {
				site.IsCurrent = true
				site.IsInService = true
				site.HasData = true
				site.CurrentAction = cmtypes.ActionTypeNone
			}
			value := strings.TrimSpace(stringTok[1])
			loc := strings.Index(value, "-")
			if loc > 0 {
				runes := []rune(value)
				site.ID = strings.TrimSpace(string(runes[0:loc]))
				location := strings.TrimSpace(string(runes[loc+1 : len(runes)]))

				runes = []rune(location)
				start := strings.LastIndex(location, "(")
				stop := strings.LastIndex(location, ")")
				if start > 0 && stop > 0 {
					site.State = string((runes[start+1 : stop]))
					location = strings.TrimSpace(string(runes[0:start]))
					site.Location = cleanLocationString(location)
				}
			}
			siteJobs <- site
		}
		rlctypes.Ctx.Log.Infof("Finished loading locations for " + state.State)
		return count
}

func loadStates() ([]rlctypes.State, error) {
	var states []rlctypes.State
	var locationJSON = filepath.FromSlash("../../resources/locations.json")
	rlctypes.Ctx.Log.Infof("Loading data for states")
	file, err := os.Open(locationJSON)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&states); err != nil {
		return nil, err
	}

	return states, nil
}

func getObservationInfo(obs string, site *cmtypes.Site) {
	if obs == "Data Is Not Current" {
		site.IsCurrent = false
		site.IsInService = true
		site.HasData = true
		site.CurrentAction = cmtypes.ActionTypeNone
	} else if obs == "No Observation Data Currently Available" {
		site.IsCurrent = false
		site.IsInService = true
		site.HasData = false
		site.CurrentAction = cmtypes.ActionTypeNone
	} else if obs == "Out of Service" {		
		site.IsCurrent = false
		site.IsInService = false
		site.HasData = false
		site.CurrentAction = cmtypes.ActionTypeNone
	} else {
		site.IsCurrent = true
		site.IsInService = true
		site.HasData = true
		if strings.Contains(obs, string(cmtypes.ActionTypeAction)){ 	
			site.CurrentAction = cmtypes.ActionTypeAction
		} else if strings.Contains(obs, string(cmtypes.ActionTypeMinor)){ 	
			site.CurrentAction = cmtypes.ActionTypeMinor
		} else if strings.Contains(obs, string(cmtypes.ActionTypeModerate)){ 	
			site.CurrentAction = cmtypes.ActionTypeModerate
		} else if strings.Contains(obs, string(cmtypes.ActionTypeMajor)){ 	
			site.CurrentAction = cmtypes.ActionTypeMajor
		} else if strings.Contains(obs, string(cmtypes.ActionTypeLowWater)){ 	
			site.CurrentAction = cmtypes.ActionTypeLowWater
		} else {
			//Unknown observation...rlctypes.Ctx.Logging it
			rlctypes.Ctx.Log.Warnf("Unknown Observation: %s", obs)
		}
	} 

}

func cleanLocationString(location string) string {
	cleanLocation := location
	for strings.Contains(cleanLocation, "(") {
		start := strings.Index(location, "(")
		stop := strings.Index(location, ")")
		if start > 0 && stop > 0 {
			runes := []rune(cleanLocation)
			startString := strings.TrimSpace(string(runes[0 : start]))
			var endString string
			if stop+1 >= len(runes) {
				endString = ""
			} else {
				endString = strings.TrimSpace(string(runes[stop+1 : len(runes)]))
			}
			cleanLocation = fmt.Sprintf("%s %s", startString, endString)
		}
	}
	return cleanLocation
}