NAME=replay_uploader
BUILDDIR=build
VERSION ?=Unknown
BuildTime:=$(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:=$(shell git rev-parse HEAD)
GOVERSION:=$(shell go version)

LDFLAGS=-w -s

GOLDFLAGS=-X 'github.com/jumpserver/replay_uploader/cmd.Version=$(VERSION)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.BuildStamp=$(BuildTime)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.GitHash=$(COMMIT)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.GoVersion=$(GOVERSION)'

GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "$(GOLDFLAGS) ${LDFLAGS}"

define make_artifact_full
	GOOS=$(1) GOARCH=$(2) $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$(1)-$(2)
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)/
	cp $(BUILDDIR)/$(NAME)-$(1)-$(2) $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$(1)-$(2).tar.gz $(NAME)-$(VERSION)-$(1)-$(2)/
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)/ $(BUILDDIR)/$(NAME)-$(1)-$(2)
endef

build:
	$(GOBUILD) -o $(BUILDDIR)/$(NAME) .

all:
	$(call make_artifact_full,darwin,amd64)
	$(call make_artifact_full,darwin,arm64)
	$(call make_artifact_full,linux,amd64)
	$(call make_artifact_full,linux,arm64)
	$(call make_artifact_full,linux,ppc64le)
	$(call make_artifact_full,linux,s390x)
	$(call make_artifact_full,linux,riscv64)

local:
	$(call make_artifact_full,$(shell go env GOOS),$(shell go env GOARCH))

darwin-amd64:
	$(call make_artifact_full,darwin,amd64)

darwin-arm64:
	$(call make_artifact_full,darwin,arm64)

linux-amd64:
	$(call make_artifact_full,linux,amd64)

linux-arm64:
	$(call make_artifact_full,linux,arm64)

linux-loong64:
	$(call make_artifact_full,linux,loong64)

linux-ppc64le:
	$(call make_artifact_full,linux,ppc64le)

linux-s390x:
	$(call make_artifact_full,linux,s390x)

linux-riscv64:
	$(call make_artifact_full,linux,riscv64)

clean:
	rm -rf $(BUILDDIR)
