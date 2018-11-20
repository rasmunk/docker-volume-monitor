.PHONY: help

OWNER:=rasmunk
TAG:=edge
IMAGE:=docker-volume-monitor

build:
	docker build --rm --force-rm -t $(OWNER)/$(IMAGE):$(TAG) .

push:
	docker push $(OWNER)/$(IMAGE):$(TAG)
