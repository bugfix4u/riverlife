// River Life
// Copyright (C) 2020  Denny Chambers

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	rlatypes "riverlife/internal/rlapisvr/types"
	appctx "riverlife/internal/rlapisvr/appcontext"
)


func main() {
 	rlatypes.Ctx = appctx.New()
	defer rlatypes.Ctx.DB.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	restSvr := rlatypes.Ctx.Run()
	<-stop
	rlatypes.Ctx.Log.Info("Shutting down the River Life API server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	restSvr.Shutdown(ctx)
	rlatypes.Ctx.Log.Info("River Life server gracefully stopped")
}
