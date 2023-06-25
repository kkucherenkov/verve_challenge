## The "feat/ms-approach" branch
The "feat/ms-approach" branch contains a microservice solution. It contains two microservices: the "storage" and "webapi". The storage service provides API to manage in-memory storage. The "webapi" service provides HTTP API according to the task description and uses the "storage" service to get data. The communication between services uses the "gRPC" protocol. This branch also contains a setup for a simple load balancer pattern via the Nginx server. Also, both services provide metrics to Prometheus.

In this solution I've used
- gin for http server
- go-kit framework for microservices
- protobuf and gRPC for communication between services
- NGinX for load balancer
- Prometheus for API monitoring
- Grafana for API monitoring visualisation

There are two docker-compose files: the load balancer version and the metrics version

## How to use and run:
1. clone the repository branch
```
git clone --branch feat/ms-approach https://github.com/kkucherenkov/verve_challenge.git
cd verve_challenge
```
2. Run the services:
```
docker-compose up --build --force-recreate
```
3. Run the test:
```
sh ./run_test.sh
```
or use curl to acces data:
```
curl --request GET \
  --url http://127.0.0.1:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
```
browser:
```
http://127.0.0.1:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
```
4. You can access the Prometheus instance via `localhost:9090`
5. You can access the Grafana instance via `localhost:3000`, use `admin/password` for credentials. Add prometheus datasource using `prometheus:9090` URL.

Test results:
```
Running 15s test @ http://webapi:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
  2 threads and 2 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    30.23ms  134.18ms 981.91ms   95.11%
    Req/Sec   788.27    189.75     1.12k    75.89%
  22310 requests in 15.01s, 4.64MB read
  Socket errors: connect 0, read 0, write 0, timeout 1
Requests/sec:   1486.26
Transfer/sec:    316.41KB
```

Use `/promotions/reload` endpoint to reload DB
```
curl --request POST \
  --url http://127.0.0.1:1321/promotions/reload \
  --header 'Content-Type: application/json' \
  --data '{
	"path": "data/promotions.csv"
}'
```
