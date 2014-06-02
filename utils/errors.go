// Copyright (c) 2014, Mauro Toffanin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
)

// (WIP) This function is an ugly and hackish error handler that needs to be
// replaced with something more complete and useful. Don't use it.
func Check(e error) {
	// TODO: replace Check() with github.com/juju/errgo helpers.
	if e != nil {
		fmt.Errorf("%v", e)
		panic(e)
	}
}
