#!/usr/bin/env node

/* Copyright 2016 Huitse Tai. All rights reserved.
 * Use of this source code is governed by BSD 3-clause
 * license that can be found in the LICENSE file.
 */

var fs = require('fs')

var modFiles = process.argv.slice(2)
var modCount = modFiles.length

var allArgs = {}

modFiles.forEach(function (val, ind, arr) {
  fs.readFile(val, 'utf8', function (err, data) {
    if (err) throw err

    var modArgs = JSON.parse(data)

    for (var prop in modArgs) {
      if (modArgs.hasOwnProperty(prop))
        allArgs[prop] = modArgs[prop]
    }

    if (--modCount <= 0) {
      process.stdout.write(JSON.stringify(allArgs, null, '\t'))
      process.stdout.write('\n')
    }
  })
})
