_@=@
BASE=$(shell pwd)
MAIN=$(BASE)/cmd
SRC_FILES=$(wildcard *.go) $(wildcard */*.go) $(wildcard */*/*.go)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitCommit=$(GIT_COMMIT)"

VERSION=$(shell $(BASE)/.tools/git-revision.sh)
BUILD_DATE=$(shell date -u +%Y/%m/%d-%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --verify HEAD)
GIT_UNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)

$(BASE)/bin/srv: $(SRC_FILES)
	$(_@) go build $(LDFLAGS) -o $@ $(MAIN)

build: $(BASE)/bin/srv

install: $(BASE)/bin/srv
	$(_@) cp $? ~/bin/

clean:
	$(_@) (cd $(BASE)/bin && ls | grep -v .gitkeep | xargs rm -rf &> /dev/null && cd $(BASE) && rm -r vendor &> /dev/null) || exit 0
