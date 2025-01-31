package proc_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/go-delve/delve/pkg/proc"
	protest "github.com/go-delve/delve/pkg/proc/test"
)

func TestGoroutineCreationLocation(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		t.Skip("test is not valid on ARM64")
	}
	protest.AllowRecording(t)
	withTestProcess("goroutinestackprog", t, func(p proc.Process, fixture protest.Fixture) {
		bp := setFunctionBreakpoint(p, t, "main.agoroutine")
		assertNoError(proc.Continue(p), t, "Continue()")

		gs, _, err := proc.GoroutinesInfo(p, 0, 0)
		assertNoError(err, t, "GoroutinesInfo")

		for _, g := range gs {
			currentLocation := g.UserCurrent()
			currentFn := currentLocation.Fn
			if currentFn != nil && currentFn.BaseName() == "agoroutine" {
				createdLocation := g.Go()
				if createdLocation.Fn == nil {
					t.Fatalf("goroutine creation function is nil")
				}
				if createdLocation.Fn.BaseName() != "main" {
					t.Fatalf("goroutine creation function has wrong name: %s", createdLocation.Fn.BaseName())
				}
				if filepath.Base(createdLocation.File) != "goroutinestackprog.go" {
					t.Fatalf("goroutine creation file incorrect: %s", filepath.Base(createdLocation.File))
				}
				if createdLocation.Line != 23 {
					t.Fatalf("goroutine creation line incorrect: %v", createdLocation.Line)
				}
			}

		}

		p.ClearBreakpoint(bp.Addr)
		proc.Continue(p)
	})
}
