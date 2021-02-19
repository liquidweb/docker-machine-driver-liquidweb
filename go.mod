module github.com/liquidweb/docker-machine-driver-liquidweb

go 1.15

require (
	github.com/docker/docker v20.10.2+incompatible // indirect
	github.com/docker/machine v0.16.2
	github.com/liquidweb/liquidweb-go v1.6.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
)

replace github.com/liquidweb/liquidweb-go => /home/ssullivan/golang/src/github.com/liquidweb/liquidweb-go
