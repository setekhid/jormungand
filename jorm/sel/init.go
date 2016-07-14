// Copyright 2016 Huitse Tai. All rights reserved.
// Use of this source code is governed by BSD 3-clause
// license that can be found in the LICENSE file.

package sel

import (
	"github.com/setekhid/jormungand/jorm"
	"github.com/setekhid/jormungand/misc/jargs"
)

func init() {

	jargs.RegistEntry(moduleName, func() { go jorm.Routing(Router()) })
}
