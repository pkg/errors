PKGS := github.com/pkg/errors
GO := go

test:
	$(GO) test $(PKGS)
