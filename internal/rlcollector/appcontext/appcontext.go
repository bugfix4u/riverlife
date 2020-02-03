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
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"os"
	dbh "riverlife/internal/common/dbhandler"
	cmtypes "riverlife/internal/common/types"
	"time"
)

const Prefix = "COLLECTOR"

type AppContext struct {
	DB     *gorm.DB
	Log    *log.Logger
	Config *Config
}

type Config struct {
	LogFormatter             string
	LogLevel                 string
	LogOutput                string
	DbUser                   string        `required:"true"`
	DbPassword               string        `required:"true"`
	DbHost                   string        `required:"true"`
	DbPort                   string        `default:"5432"`
	DbName                   string        `required:"true"`
	StateTickerTimeHour      time.Duration `default:"24"`
	SiteTickerTimeHour       time.Duration `default:"1"`
	StateChannelSize         int32         `default:"60"`
	SiteChannelSize          int32         `default:"10000"`
	PersistChannelSize       int32         `default:"10000"`
	StateWorkerThreadCount   int           `default:"10"`
	SiteWorkerThreadCount    int           `default:"10"`
	PersistWorkerThreadCount int           `default:"10"`
	RedisHost                string        `required:"true"`
	RedisPort                string        `required:"true"`
}

func New() *AppContext {
	var newCtx AppContext
	newCtx.initializeConfig()
	newCtx.initializeLogger()
	newCtx.initializeDB()
	newCtx.Log.WithFields(log.Fields{
		"LogFormatter":             newCtx.Config.LogFormatter,
		"LogLevel":                 newCtx.Config.LogLevel,
		"LogOutput":                newCtx.Config.LogOutput,
		"DbUser":                   newCtx.Config.DbUser,
		"DbPassword":               "*********",
		"DbHost":                   newCtx.Config.DbHost,
		"DbPort":                   newCtx.Config.DbPort,
		"DbName":                   newCtx.Config.DbName,
		"StateTickerTimeHour":      newCtx.Config.StateTickerTimeHour,
		"SiteTickerTimeHour":       newCtx.Config.SiteTickerTimeHour,
		"StateChannelSize":         newCtx.Config.StateChannelSize,
		"SiteChannelSize":          newCtx.Config.SiteChannelSize,
		"PersistChannelSize":       newCtx.Config.PersistChannelSize,
		"StateWorkerThreadCount":   newCtx.Config.StateWorkerThreadCount,
		"SiteWorkerThreadCount":    newCtx.Config.SiteWorkerThreadCount,
		"PersistWorkerThreadCount": newCtx.Config.PersistWorkerThreadCount,
		"RedisHost":                newCtx.Config.RedisHost,
		"RedisPort":                newCtx.Config.RedisPort,
	}).Debug("Collector Configuration Settings")
	return &newCtx
}

func (asc *AppContext) initializeConfig() {
	var conf Config
	err := envconfig.Process(Prefix, &conf)
	if err != nil {
		log.Fatal(err.Error())
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
	var err error
	asc.DB = dbh.GetDbConnection(asc.Config.DbUser,
		asc.Config.DbPassword,
		asc.Config.DbName,
		asc.Config.DbHost,
		asc.Config.DbPort)

	asc.DB.SetLogger(asc.Log)

	if asc.Config.LogLevel == "debug" || asc.Config.LogLevel == "trace" {
		asc.DB.LogMode(true)
	}

	if err != nil {
		log.Fatal(err)
	}
	asc.DB.SingularTable(true)
	asc.DB.AutoMigrate(
		&cmtypes.Site{},
	)

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
