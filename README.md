# Wormhole Wallet API Service



## What is Wormhole

Wormhole is a fast, portable Omni Layer implementation that is based off the Bitcoin ABC codebase (currently 0.18.2). This implementation requires no external dependencies extraneous to Bitcoin ABC, and is native to the Bitcoin Cash network just like other Bitcoin Cash nodes. It currently supports a wallet mode and is seamlessly available on three platforms: Windows, Linux and Mac OS. Wormhole Cash Layer extensions are exposed via the JSON-RPC interface. Development has been consolidated on the Wormhole product, and it is the reference client for the Wormhole Cash Layer.

## Quick Start

#### Prerequesites

| Package | Version |
| :------ | ------- |
| Mysql   |         |
| Golang  |         |
| Redis   |         |

#### Database Init

1、Once the mysql installation is complete.you need create your database first.

```
#connect your database
mysql -u{your-user} -p{your-password} -h{your-host} -P{your-port}

#create database
mysql> create database wormhole;
```

2、Init tables

```
git clone https://github.com/copernet/whccommon
cd model/operation/setup/
go build
./setup --host={your-host} --port={your-port} --user={your-user} --passwd={your-password} --db=wormhole
```

#### Config Init

```
git clone https://github.com/copernet/whcwallet
cp conf.yml.sample conf.yml

#you need modify db、redis、rpc、log to your local config
```

#### How To Run

```
cd ${gopath}/src/github.com/copernet/whcwallet
mkdir logs
go build

#start
tools/run start whcwallet
#stop 
tools/run stop whcwallet
```

## Document
[whcwallet](https://github.com/EmbraceDD/whcwallet/blob/master/doc/whcwallet.md)

