## The "feat/redis" branch
The "feat/redis" branch contains a solution that uses external storage, Redis in this case. There is an API service that provides HTTP API and has no state, so it allows horizontal scaling, but it also involves network communication overhead. This branch also contains a setup for a simple load balancer pattern via the Nginx server.

## How to use and run:
1. clone the repository branch
```
git clone --branch feat/redis https://github.com/kkucherenkov/verve_challenge.git
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
Results for singe instance:
```
Running 15s test @ http://ngnix-server:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
  4 threads and 4 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    32.50ms  139.96ms   1.01s    94.93%
    Req/Sec   552.93    129.76   818.00     72.34%
  31276 requests in 15.02s, 8.59MB read
Requests/sec:   2081.94
Transfer/sec:    585.54KB
```
Results for load balancer version:
```
Running 15s test @ http://app1:8080/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
  4 threads and 4 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   767.96us    1.40ms  41.04ms   98.20%
    Req/Sec     1.55k   299.89     2.85k    73.83%
  92394 requests in 15.00s, 21.32MB read
  Socket errors: connect 0, read 0, write 0, timeout 4
Requests/sec:   6158.05
Transfer/sec:      1.42MB
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
