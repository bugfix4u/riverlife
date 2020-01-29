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
package dbhandler

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Gorm wants to know the type of database
var dialect = "postgres"

func GetDbConnection(dbUser, dbPassword, dbName, dbHost, dbPort string) *gorm.DB {
	connectString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port= %s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	conn, err := gorm.Open(dialect, connectString)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func DbResults(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
