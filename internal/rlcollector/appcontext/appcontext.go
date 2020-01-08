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

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	dbh "riverlife/internal/common/dbhandler"
	cmtypes "riverlife/internal/common/types"
)

type AppContext struct {
	DB        *gorm.DB
	Log       *log.Logger
	LogOutput string
}

func New() *AppContext {
	var newCtx AppContext
	newCtx.initializeFlags()
	newCtx.initializeLogger()
	newCtx.initializeDB()
	return &newCtx
}

func (asc *AppContext) initializeFlags() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp"
	}
	asc.LogOutput = fmt.Sprintf("%s/rlcollector.log", home)
}

func (asc *AppContext) initializeLogger() {

	asc.Log = log.New()
	asc.Log.Formatter = new(log.JSONFormatter)
	asc.Log.Level = log.TraceLevel
	if asc.LogOutput != "" {
		file, err := os.OpenFile(asc.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			asc.Log.Out = file
		} else {
			asc.Log.Out = os.Stdout
			asc.Log.Info("Failed to log to file, using default stdout")
		}
	} else {
		asc.Log.Out = os.Stdout
	}
}

func (asc *AppContext) initializeDB() {
	var err error
	asc.DB = dbh.GetDbConnection()

	if err != nil {
		log.Fatal(err)
	}
	asc.DB.SingularTable(true)
	asc.DB.Debug().AutoMigrate(
		&cmtypes.Site{},
	)

}
