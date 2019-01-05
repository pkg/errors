PKGS := github.com/pkg/errors
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt unused misspell unconvert gosimple ineffassign

test: 
	$(GO) test $(PKGS)

vet: | test
	$(GO) vet $(PKGS)

staticcheck:
	$(GO) get honnef.co/go/tools/cmd/staticcheck
	staticcheck $(PKGS)

unused:
	$(GO) get honnef.co/go/tools/cmd/unused
	unused -exported $(PKGS)

misspell:
	$(GO) get github.com/client9/misspell/cmd/misspell
	misspell \
		-locale GB \
		-error \
		*.md *.go

unconvert:
	$(GO) get github.com/mdempsky/unconvert
	unconvert -v $(PKGS)

gosimple:
	$(GO) get honnef.co/go/tools/cmd/gosimple
	gosimple $(PKGS)

ineffassign:
	$(GO) get github.com/gordonklaus/ineffassign
	find $(SRCDIRS) -name '*.go' | xargs ineffassign

pedantic: check unparam errcheck staticcheck

unparam:
	$(GO) get mvdan.cc/unparam
	unparam ./...

errcheck:
	$(GO) get github.com/kisielk/errcheck
	errcheck $(PKGS)

gofmt:  
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"
