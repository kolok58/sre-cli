# Manage Kubernetes deployments using kolok58/sre-cli

## Installation

```sh
git clone https://github.com/kolok58/sre-cli.git
```

## Build command
Ensure you are within the code directory
```
go build -o sre 
```

## Usage

This tool assumes you have already authenticated to Kubernetes and have a valid .kube/config file
```
./sre
```

Alternatively you can move the binary to a directory in your PATH and call the binary directly:
 
```
mv /path/to/repo /usr/local/bin
sre -h
```