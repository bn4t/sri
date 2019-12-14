SHELL = /bin/sh
.SUFFIXES:
.SUFFIXES: .go

PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin
GO?=go
GOFLAGS?=

GOSRC!=find . -name '*.go'


all: build


# Exists in GNUMake but not in NetBSD make and others.
RM?=rm -f

build: $(GOSRC)
	$(GO) build $(GOFLAGS) -o sri

install: build
	install -m755 sri $(BINDIR)/sri


uninstall:
	$(RM) $(BINDIR)/sri

clean: 
	$(RM) sri

.DEFAULT_GOAL := all
.PHONY: all build install uninstall clean
