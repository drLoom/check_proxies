# It is a simle tool, to check proxies against test url.

Proxies to check must be stored in a CSV file like:
```
head -n5 ./data/test_prxoxies.csv
```

| IP Address    | Port|
| ------------- |-----|
|93.64.183.162|56508|
|92.220.168.252|80|
|82.77.114.201|3128|
|52.233.190.21|80|
...

Output file is going to be like:
```
head -n5 ./data/checked_proxies.csv 
```
| IP|Port|HTTP Code|Response Duration|Error message|
| ------------- |-----| ------------- |-----|-----|
|81.24.82.69|37016|0|0.055434|Head http://www.profinance.ru/: proxyconnect tcp: dial tcp 81.24.82.69:37016: connect: no route to host|
|91.203.114.105|45770|0|0.055207|Head http://www.profinance.ru/: proxyconnect tcp: dial tcp 91.203.114.105:45770: connect: connection refused|
|82.100.4.63|30595|0|0.092923|Head http://www.profinance.ru/: proxyconnect tcp: dial tcp 82.100.4.63:30595: connect: connection refused|
|40.91.214.91|80|405|0.127598|

```
go run main.go -h
Usage of /tmp/go-build170132499/b001/exe/main:
  -output string
    	Folder, result file will be stored to
  -source string
    	Path to CSV file with proxies to check
  -threads int
    	Number of threads (default 5)
  -timeout int
    	Timeout per request, s (default 5)
  -url string
    	URL to test proxy
  -verbose
    	Log output
```
Command line example:
```
go run main.go \
    -source=./data/test_prxoxies.csv \
    -threads=50 \
    -url=http://www.profinance.ru/ \
    -timeout=5 \
    -output=./data/ \
    -verbose=true
```

Also, can be run in docker (see Dockerfile):
```
docker run \
    --rm \
    -u $(id -u):$(id -g) \
    -v $(find "$(pwd)" -name 'data'):/data/ \
    check_proxies \
        -source=/data/test_prxoxies.csv \
        -threads=250 \
        -url=http://www.profinance.ru \
        -timeout=10 \
        -output=/data \
        -verbose=true
```