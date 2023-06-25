#!/bin/sh

docker run --network=backend \
        --rm skandyla/wrk -t2 -c2 -d90s http://webapi:1321/promotions/33ae574d-60d3-4d65-9deb-f25fd04a229a
