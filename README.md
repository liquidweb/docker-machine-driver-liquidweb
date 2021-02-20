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

### Set Liquid Web API Credentials

```shell
export LW_USERNAME='username'
export LW_PASSWORD='password'
```

These alternatively can be given as an argument to `docker-machine` (`--lw-username` and `--lw-password`)

### Starting/Stopping/Restarting machines

docker-machine allows you to perform basic actions on a compute resource. For example starting, stopping, and restarting is all supported. See `docker-machine help` for more options.

```shell
docker-machine {restart|start|stop} docker01
```

### Spinning up a new compute node

Create a new compute node with template `UBUNTU_2004_UNMANAGED` (overriding the default)

```shell
user@host $ docker-machine create docker04 --driver liquidweb --lw-template UBUNTU_2004_UNMANAGED --lw-zone-id 12
Running pre-create checks...
(docker04) Generating a random root password...
Creating machine...
(docker04) Creating liquidweb machine instance...
(docker04) Created liquidweb instance with uniq_id [KRGV3V]
(docker04) Waiting for machine to become ready..
(docker04) node [KRGV3V] has become ready
(docker04) Machine has become ready..
(docker04) Discovered IP address [50.28.32.147] for uniq_id [KRGV3V]
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Error creating machine: Error running provisioning: Unable to verify the Docker daemon is listening: Maximum number of retries (10) exceeded
user@host $ 
```

If you get an error like above, you can retry the docker installation:

```shell
user@host $ docker-machine provision docker04
Waiting for SSH to be available...
Detecting the provisioner...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
user@host $ 
```

Once successful, you should see the node status as `READY`:

```shell
user@host $ docker-machine ls
NAME       ACTIVE   DRIVER      STATE     URL                       SWARM   DOCKER     ERRORS
docker01   -        liquidweb   Running   tcp://50.28.52.211:2376           v20.10.3   
docker04   -    liquidweb   Running   tcp://50.28.32.147:2376        v20.10.3   
user@host $ 
```

### Deleting a compute node

```shell
user@host $ docker-machine rm docker05
About to remove docker05
WARNING: This action will delete both local reference and remote instance.
Are you sure? (y/n): y
(docker05) removing liquidweb compute node [DTHPDU]...
Successfully removed docker05
user@host $ 
```

### Pointing to your machine

Configure your local docker client to use your remote machine:

```shell
user@host $ eval $(docker-machine env docker01)
user@host $ docker-machine ls
NAME       ACTIVE   DRIVER      STATE     URL                      SWARM   DOCKER     ERRORS
docker01   *        liquidweb   Running   tcp://50.28.52.94:2376           v20.10.3   
user@host $ 
```
Note that in the `docker-machine ls` above, there is an `*` to the right of `docker01` because its the active host from the eval above. Any containers created/started via your local `docker` will now happen on the remote `docker01`, and not localhost.
