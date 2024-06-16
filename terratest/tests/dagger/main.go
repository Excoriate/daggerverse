package main

import "context"

type Tests struct {
	ModSrc *Directory
}

func New(
	// src is the directory that contains all the source code, including the module directory.
	src *Directory,
) *Tests {
	//if src == nil {
	//	src = dag.Module().Source().ContextDirectory()
	//}
	return &Tests{
		ModSrc: src,
	}
}

func (m *Tests) TestDaggerCall() (string, error) {
	return dag.Dagindag().
		DagCli(context.Background(), DagindagDagCliOpts{
			DagCmds: "call",
			Src:     m.ModSrc,
		})
}
