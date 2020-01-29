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
	"os"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	dbh "riverlife/internal/common/dbhandler"
	mh "riverlife/internal/rlapisvr/muxhandler"
	cmtypes "riverlife/internal/common/types"
	log "github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

const Prefix = "APISVR"

type AppContext struct {
	Router *mux.Router
	DB     *gorm.DB
	Log    *log.Logger
	Config *Config
}

type Config struct {
	LogFormatter 							string
	LogLevel 									string
	LogOutput 								string
	DbUser										string `required:"true"`
	DbPassword 								string `required:"true"`
	DbHost 										string `required:"true"`
	DbPort 										string `default:"5432"`
	DbName 										string `required:"true"`
	IsHttps										bool `default:"false"`
	HttpPort									string
	DefaultPageCount					int32 `default:"500"`
	TlsCertFile								string
	TlsKeyFile								string
}

func New() *AppContext {
	var newCtx AppContext
	newCtx.initializeConfig()
	newCtx.initializeLogger()
	newCtx.initializeDB()
	newCtx.initializeRouter()
	newCtx.Log.WithFields(log.Fields{
		"LogFormatter": newCtx.Config.LogFormatter,
		"LogLevel": newCtx.Config.LogLevel,
		"LogOutput": newCtx.Config.LogOutput,
		"DbUser": newCtx.Config.DbUser,
		"DbPassword": "*********",
		"DbHost": newCtx.Config.DbHost,
		"DbPort": newCtx.Config.DbPort,
		"IsHttps": newCtx.Config.IsHttps,
		"HttpPort": newCtx.Config.HttpPort,
		"DefaultPageCount": newCtx.Config.DefaultPageCount,
		"TlsCertFile": newCtx.Config.TlsCertFile,
		"TlsKeyFile": newCtx.Config.TlsKeyFile,
	}).Debug("Apisvr Configuration Settings")
	return &newCtx
}

func (asc *AppContext) initializeConfig() {
	var conf Config
	err := envconfig.Process(Prefix, &conf)
  if err != nil {
  	log.Fatal(err.Error())
	}

	if conf.IsHttps {
		if conf.TlsCertFile == "" || conf.TlsKeyFile == "" {
			log.Fatal("To enable HTTPS, you must provide a certificate file and a key file.")
		}
	}

	if conf.HttpPort == "" {
		if conf.IsHttps {
			conf.HttpPort = ":8443"
		} else {
			conf.HttpPort = ":8080"
		}
	}

	asc.Config = &conf
	
}

func (asc *AppContext) initializeLogger() {
	asc.Log = log.New()
	asc.Log.SetFormatter(asc.Config.getLogFormatter())
	level, err := log.ParseLevel(asc.Config.LogLevel)
	if err != nil {
		log.Fatal(err.Error())
	}
	asc.Log.SetLevel(level)
	if asc.Config.LogOutput != "" {
		file, err := os.OpenFile(asc.Config.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	asc.Log.Info("Setting up DB connection")
	asc.DB = dbh.GetDbConnection(	asc.Config.DbUser, 
																asc.Config.DbPassword, 
																asc.Config.DbName, 
																asc.Config.DbHost, 
																asc.Config.DbPort,
	)
}

func (asc *AppContext) initializeRouter() {
	asc.Router = mux.NewRouter()

	sm := mh.NewSiteMux(asc.DB, asc.Log)
	asc.Log.Infof("Setting up mux handlers for %s", sm.GetName())
	sm.InitRouter(asc.Router)

}

func (asc *AppContext) Run() *http.Server {
	asc.Log.Info("Starting the River Life API server on port" + asc.Config.HttpPort)
	svr := &http.Server{Addr: asc.Config.HttpPort, Handler: asc.Router}

	go func() {
		if asc.Config.IsHttps {
			log.Fatal(svr.ListenAndServeTLS(asc.Config.TlsCertFile, asc.Config.TlsKeyFile))
		} else {
			log.Fatal(svr.ListenAndServe())
		}
	}()

	return svr

}

func (conf *Config) getLogFormatter() log.Formatter {
	switch conf.LogFormatter {
	case cmtypes.Json:
		return new(log.JSONFormatter)
	case cmtypes.Text:
		return new(log.TextFormatter)
	default:
		return new(log.JSONFormatter)
	}
}
