.PHONY: help

OWNER:=ucphhpc
TAG:=edge
IMAGE:=docker-volume-monitor
ARGS=

build:
	docker build --rm --force-rm -t $(OWNER)/$(IMAGE):$(TAG) $(ARGS) .

push:
	docker push $(OWNER)/$(IMAGE):$(TAG) $(ARGS)
