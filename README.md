# uscan
`universal blockchain scan for EVM series`


### 一、 binary & running
```
make all-build
bin/uscan --rpc_urls "wss://testnet.ankr.com/ws" 
```

### 二、docker image & running
```
make docker-build
docker run -it -p 4322:4322 uchainorg/uscan:latest --rpc_urls "wss://testnet.ankr.com/ws" 
```



