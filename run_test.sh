#!/bin/sh

docker run --network=backend \
        --rm skandyla/wrk -t4 -c4 -d15s http://app:8080/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
