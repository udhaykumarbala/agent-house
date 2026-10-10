package web

import (
	"testing"
	"time"
)

// Chat-initiated builds used to construct a private orchestrator with no
// checkpoint wiring and no agent registry, so semi_auto gates hung. The fix
// routes the Brain's create_project through the same launchBuild path as
// /api/task. This test pins the delegation.
func TestBrainCreateProjectDelegatesToLaunchBuild(t *testing.T) {
	got := make(chan [2]string, 1)
	bh := &BrainHandler{
		projDir: t.TempDir(),
		LaunchBuild: func(projectID, task, parentTaskID, taskID string) {
			got <- [2]string{projectID, task}
		},
	}
	if err := bh.handleCreateProject("fitflow-app", "Build the FitFlow MVP"); err != nil {
		t.Fatal(err)
	}
	select {
	case v := <-got:
		if v[0] != "fitflow-app" || v[1] != "Build the FitFlow MVP" {
			t.Fatalf("wrong delegation: %v", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("LaunchBuild was never called — create_project still uses the unwired private orchestrator")
	}
}
