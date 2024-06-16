package main

import (
	"context"
)

type Tests struct {
	TestDir *Directory
}

func New() *Tests {
	t := &Tests{}
	t.TestDir = t.getTestDir()
	return t
}

// getTestDir returns the test directory.
//
// This is a helper method for tests, in order to get the test directory which
// is located in the same directory as the test file, and normally named as "testdata".
func (m *Tests) getTestDir() *Directory {
	return dag.CurrentModule().Source().Directory("./testdata")
}

func (m *Tests) TestModuleCall() (string, error) {
	return dag.Dagindag().DagCli(context.Background(),
		DagindagDagCliOpts{
			DagCmds: "call",
			Src:     m.TestDir,
		})
}
