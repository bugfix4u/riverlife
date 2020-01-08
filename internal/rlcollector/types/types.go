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
	"encoding/xml"
	"sync"
)

type State struct {
	State string
	Abbr  string
}

type NoaaSite struct {
	XMLName        xml.Name `xml:"site"`
	Text           string   `xml:",chardata"`
	Xsi            string   `xml:"xsi,attr"`
	Timezone       string   `xml:"timezone,attr"`
	Originator     string   `xml:"originator,attr"`
	Name           string   `xml:"name,attr"`
	ID             string   `xml:"id,attr"`
	Generationtime string   `xml:"generationtime,attr"`
	Sigstages      Sigstage `xml:"sigstages"`
	Sigflows       Sigflow  `xml:"sigflows"`
	Observed       Observed `xml:"observed"`
	Forecast       Forecast `xml:"forecast"`
}

type Sigstage struct {
	Text     string `xml:",chardata"`
	Low      Stage  `xml:"low"`
	Action   Stage  `xml:"action"`
	Bankfull Stage  `xml:"bankfull"`
	Flood    Stage  `xml:"flood"`
	Moderate Stage  `xml:"moderate"`
	Major    Stage  `xml:"major"`
	Record   Stage  `xml:"record"`
}

type Stage struct {
	Text  string `xml:",chardata"`
	Units string `xml:"units,attr"`
}

type Sigflow Sigstage

type Observed struct {
	Text   string  `xml:",chardata"`
	Datums []Datum `xml:"datum"`
}

type Sample struct {
	Text  string `xml:",chardata"`
	Name  string `xml:"name,attr"`
	Units string `xml:"units,attr"`
}

type Valid struct {
	Text     string `xml:",chardata"`
	Timezone string `xml:"timezone,attr"`
}

type Datum struct {
	Text      string `xml:",chardata"`
	Valid     Valid  `xml:"valid"`
	Primary   Sample `xml:"primary"`
	Secondary Sample `xml:"secondary"`
	Pedts     string `xml:"pedts"`
}

type Forecast struct {
	Text     string  `xml:",chardata"`
	Timezone string  `xml:"timezone,attr"`
	Issued   string  `xml:"issued,attr"`
	Datums   []Datum `xml:"datum"`
}

type SafeCount struct {
	Count int64
	mux   sync.Mutex
}

func (sc *SafeCount) IncCount(count int64) {
	sc.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	sc.Count += count
	sc.mux.Unlock()
}
