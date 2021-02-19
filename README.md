# docker-machine-driver-liquidweb
Liquid Web Docker machine driver

## Installation

You can grab the driver precompiled for various platforms on this projects [releases page](https://github.com/liquidweb/docker-machine-driver-liquidweb/releases).

## Installing from source

To install from source code run `make` in the root of this repository. This will install `docker-machine-driver-liquidweb` into your `${GOPATH}/bin`.

## Why use this?

If you are a [Liquid Web](https://liquidweb.com) customer with a Cloud VPS who utilizes [docker](https://docker.com) or [docker-machine](https://docs.docker.com/machine/) this tool allows you to easily spin up Liquid Web Cloud VPS servers for the compute nodes to run your containers. Rather than having to use the [Liquid Web API](https://cart.liquidweb.com/storm/api/docs/bleed), [liquidweb-cli](https://github.com/liquidweb/liquidweb-cli), or our management portal to create servers, you can use the familar framework of docker-machine. Docker-machine makes it easier to manage multiple docker hosts from your local docker client. 

## Usage

This is a [docker-machine](https://docs.docker.com/machine/) driver, so make sure you have docker and docker-machine installed as a prerequisite.

### Listing machines

```shell
user@host $ docker-machine ls
NAME       ACTIVE   DRIVER      STATE     URL                        SWARM   DOCKER     ERRORS
docker01   -        liquidweb   Running   tcp://69.167.152.19:2376           v20.10.2   
docker02   -        liquidweb   Running   tcp://209.59.138.4:2376            v20.10.2   
docker03   -        liquidweb   Running   tcp://67.227.190.28:2376           v20.10.2   
docker04   -        liquidweb   Running   tcp://67.225.160.31:2376           v20.10.2   
docker05   -        liquidweb   Running   tcp://67.227.198.22:2376           v20.10.2   
docker06   -        liquidweb   Running   tcp://209.59.129.37:2376           v20.10.2   
user@host $
```

### Starting/Stopping/Restarting machines

docker-machine allows you to perform basic actions on a compute resource. Below has an example of a graceful restart.

```shell
docker-machine {restart|start|stop} docker01
```
