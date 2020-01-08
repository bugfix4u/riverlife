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

package appcontext

import  (
	"log"
	"net/http"
	mh "riverlife/internal/rlapisvr/muxhandler"
	dbh "riverlife/internal/common/dbhandler"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type AppContext struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (asc *AppContext) InitializeDB() {
	asc.DB = dbh.GetDbConnection()
}

func (asc *AppContext) InitializeRouter() {
	asc.Router = mux.NewRouter()

	sm := mh.NewSiteMux(asc.DB)
	log.Printf("Setting up mux handlers for %s\n", sm.GetName())
	sm.InitRouter(asc.Router)
	
}

func (asc *AppContext) Run(port string) *http.Server {
	svr := &http.Server{Addr: port, Handler: asc.Router}

	go func() {
		log.Fatal(svr.ListenAndServe())
	}()

	return svr
}