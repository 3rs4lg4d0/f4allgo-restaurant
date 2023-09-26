#!/bin/bash

# -----------------------------------------------------------------------------
# This benchmark registers new restaurants.
# -----------------------------------------------------------------------------
ab -n 100000 -c 40 -p postRestaurant.json -T 'application/json' \
-H 'Content-Type: application/json' http://localhost:8080/api/v1/restaurants
