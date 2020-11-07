COMPONENTS ?= oscli

include ./dubplate.Makefile
include ./go.Makefile

snapshot:
	goreleaser --snapshot --rm-dist

release:
	goreleaser --rm-dist