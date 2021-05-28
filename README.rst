====================
docker-volume-monitor
====================

Docker image that continuously checks for docker volumes on a given host, has two options for execution:

- build and execute the binary native on the host itself.
- run as a container with the hosts docker socket mounted inside the container.

---------------
Getting Started
---------------

Either use docker to build an image, or build the binary and run it on the host.
With go example ::

    # Prune every docker volume every 10 minutes (the default)
    go build ./...
    ./docker-volume-monitor -prune-unused -interval 10
    
    # If only volumes should be monitored/listed
    ./docker-image-updater -interval 10

Build as a docker image (defaults to use the :edge tag)::

    make build
    
    # override the build tag, e.g
    make build TAG=latest

Which produces an image called nielsbohr/docker-volume-monitor:edge by default, override the TAG variable in the makefile to change this.
To run a monitor container that continuously checks and removes unused volumes ::

    docker run --mount type=bind,src=/var/run/docker.sock,target=/var/run/docker.sock nielsbohr/docker-volume-monitor:edge -prune-unused -interval 10

