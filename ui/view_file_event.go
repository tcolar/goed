package ui

import (
	"path"
	"path/filepath"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

// Handle filewatcher events
func (v *View) fileEvent(op core.FileOp, loc string) {
	if v.Backend() == nil {
		return
	}
	loc, _ = filepath.Abs(loc)
	parent := path.Dir(loc)
	wd, _ := filepath.Abs(v.workDir)
	src, _ := filepath.Abs(v.Backend().SrcLoc())
	switch v.Type() {
	case core.ViewTypeStandard:
		if src == loc && op == core.OpChmod {
			// reload, idf dirty, ask
			if !v.dirty {
				actions.Ar.ViewReload(v.id)
			}
		}
		if src == loc && (op == core.OpRemove || op == core.OpRename) {
			// close view, if dirty, ask
			if !v.dirty {
				actions.Ar.EdDelView(v.id, true)
			}
		}
	case core.ViewTypeDirListing:
		if parent == wd && (op == core.OpCreate || op == core.OpRemove || op == core.OpRename) {
			actions.Ar.ViewReload(v.id)
		}
		if loc == wd && (op == core.OpRename || op == core.OpRemove) {
			actions.Ar.EdDelView(v.id, true)
		}
	case core.ViewTypeShell:
		// nothing, maybe close the shell if CWD is gone ??
	case core.ViewTypeCmdOutput:
		// nothing
	}
}
