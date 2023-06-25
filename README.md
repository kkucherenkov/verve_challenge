# The "main" branch
The "main" branch contains a simple naive solution. It's a monolith app, that provides HTTP API according to the task description. It uses embedded in-memory storage.

Advantages:
- It's simple
- There is no network communication overhead.
  Disadvantages:
- It's not scalable

## How to use and run:
1. clone the repository branch
```
git clone https://github.com/kkucherenkov/verve_challenge.git
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
Results for singe instance:
```
Running 15s test @ http://app:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
  4 threads and 4 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.12ms    8.33ms 210.11ms   98.91%
    Req/Sec     2.63k   618.93     5.07k    74.12%
  157116 requests in 15.01s, 36.26MB read
Requests/sec:  10470.65
Transfer/sec:      2.42MB
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
