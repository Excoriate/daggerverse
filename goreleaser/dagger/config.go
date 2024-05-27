package main

const (
	goReleaserDefaultVersion = "latest"
	goReleaserDefaultImage   = "goreleaser/goreleaser"
	goReleaserDefaultCfgFile = ".goreleaser.yaml"
)

func setToDefaultCfgIfEmpty(cfgFile string) string {
	if cfgFile == "" {
		return goReleaserDefaultCfgFile
	}
	return cfgFile
}
