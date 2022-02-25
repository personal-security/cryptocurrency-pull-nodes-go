# coin-nodes-rest-go

In process dev

## Bitcoind setup

https://gist.github.com/rjmacarthy/b56497a81a6497bfabb1

## Configs files

Create all files in root folder config/*

### Variables

IP - ip address

PORT - port for this address

Need use user, pass for bitcoind

### bitcoin.json

```JSON
{
    "bitcoin_rpc_host":"IP",
    "bitcoin_rpc_port":"PORT",
    "bitcoin_rpc_user":"",
    "bitcoin_rpc_pass":""
}
```

### ethereum.json

```JSON
{
    "ethereum_rpc_host":"http://IP:PORT",
    "ethereum_rpc_user":"",
    "ethereum_rpc_pass":""
}
```

### keys.json

```JSON
{
    "api_key":""
}
```

### network.json

```JSON
{
    "ip":"IP",
    "port":"PORT"
}
```

## Run inside docker container

1. Build docker image.  
`docker build . -t coin-nodes-rest`  
2. Run docker container. _App will be launched on port number 3000._  
`docker run -t -p 8000:8000 coin-nodes-rest`
