// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package atexit

import (
	"os"
	"os/signal"
)

func StartCleanRoutine() {

	termSigs := make(chan os.Signal, 2)
	signal.Notify(termSigs, os.Interrupt)
	signal.Notify(termSigs, os.Kill)

	go func() {

		for sig := range termSigs {

			switch sig {

			case os.Interrupt:
				Exit()
				return

			case os.Kill:
				Exit()
				return
			}
		}
	}()
}
