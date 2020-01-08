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

package muxhandler

import (
	"database/sql"
	"errors"
	dbh "riverlife/internal/common/dbhandler"
	cmtypes "riverlife/internal/common/types"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type SiteMux struct {
	name string
	db   *gorm.DB
}

func NewSiteMux(sdb *gorm.DB) *SiteMux {
	site := SiteMux{
		name: "SiteMux",
		db:   sdb,
	}

	return &site
}

func (sm *SiteMux) GetName() string {
	return sm.name
}

func (sm *SiteMux) InitRouter(router *mux.Router) {
	if router == nil {
		log.Fatal(errors.New("Fatal: null router"))
	}
	router.HandleFunc("/api/v1/sites", sm.getSites).Methods("GET")
	router.HandleFunc("/api/v1/site/{id}", sm.getSite).Methods("GET")
}

func (sm *SiteMux) getSites(w http.ResponseWriter, r *http.Request) {
	Sites, err := dbh.GetSites(sm.db)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, Sites)
}

func (sm *SiteMux) getSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		RespondWithError(w, http.StatusBadRequest, "Invalid Site ID")
		return
	}

	site := cmtypes.Site{}
	site.ID = id
	if err := dbh.GetSite(sm.db, &site); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "Site not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, site)
}