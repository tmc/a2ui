package a2uibuild

import "github.com/tmc/a2ui/v09"

// Children returns a static child list containing ids.
func Children(ids ...string) v09.ChildList {
	return v09.ChildList{IDs: append([]string(nil), ids...)}
}
