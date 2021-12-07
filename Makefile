NAME=replay_uploader
BUILDDIR=build
VERSION ?=Unknown
BuildTime:=$(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:=$(shell git rev-parse HEAD)
GOVERSION:=$(shell go version)
GOLDFLAGS=-X 'github.com/jumpserver/replay_uploader/cmd.Version=$(VERSION)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.BuildStamp=$(BuildTime)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.GitHash=$(COMMIT)'
GOLDFLAGS+=-X 'github.com/jumpserver/replay_uploader/cmd.GoVersion=$(GOVERSION)'

GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "$(GOLDFLAGS)"


PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm64

WINDOWS_ARCH_LIST = \
	windows-amd64

all-arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@.exe
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@.exe $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME).exe
	cd $(BUILDDIR) && tar -czvf $(NAME)-$(VERSION)-$@.tar.gz $(NAME)-$(VERSION)-$@
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@.exe

clean:
	rm -rf $(BUILDDIR)