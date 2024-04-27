package main

import "fmt"

func (g *Goreleaser) resolveCfgArg(cfg string) string {
	var cfgFileResolved string
	if cfg != "" && cfg != goReleaserDefaultCfgFile {
		cfgFileResolved = cfg
	} else {
		if g.CfgFile != "" {
			cfgFileResolved = g.CfgFile
		} else {
			cfgFileResolved = goReleaserDefaultCfgFile
		}
	}

	return fmt.Sprintf("--config=%s", cfgFileResolved)
}
