// Package daemon implements a service for mediating access to the data store,
// and its client.
//
// Most RPCs exposed by the service correspond to the methods of Store in the
// store package and are not documented here.
package daemon

import "github.com/m9rco/phoenix-shell/src/pkg/util"

var logger = util.GetLogger("[daemon] ")

// Version is the API version. It should be bumped any time the API changes.
const Version = -93
