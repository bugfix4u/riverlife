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
	appctx "riverlife/internal/rlcollector/appcontext"
	collector "riverlife/internal/rlcollector/collector"
	rlctypes "riverlife/internal/rlcollector/types"
	"sync"
	"syscall"
)

func main() {
	rlctypes.Ctx = appctx.New()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go collector.DoStartupCollection(ctx, &wg)
	wg.Add(1)
	go collector.ScheduleStateCollection(ctx, &wg)
	wg.Add(1)
	go collector.ScheduleSiteCollection(ctx, &wg)
	<-stop
	rlctypes.Ctx.Log.Printf("Starting shutdown the River Life collection server...")
	cancel()
	wg.Wait()
	rlctypes.Ctx.Log.Printf("River Life collection server shutdown gracefully")
}
