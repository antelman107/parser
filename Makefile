GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)
BINARY=$(GOBIN)/app

build:
	@echo "  >  Building binary..."
	@GOBIN=$(GOBIN) go build -o $(BINARY) $(GOFILES)

top1:
	@echo "  >  Top 10 active users sorted by amount of PRs created and commits pushed..."
	@$(BINARY) top

top2:
	@echo "  >  Top 10 repositories sorted by amount of commits pushed..."
	@$(BINARY) top --events_entity_column_index 3 --event_types 'PushEvent' --entity_file "./data/repos.csv"

top3:
	@echo "  >  Top 10 repositories sorted by amount of watch events..."
	@$(BINARY) top --events_entity_column_index 3 --event_types 'WatchEvent' --entity_file "./data/repos.csv"
