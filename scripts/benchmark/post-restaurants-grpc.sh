#!/bin/bash

# -----------------------------------------------------------------------------
# This benchmark registers new restaurants.
# -----------------------------------------------------------------------------
k6 run -u 200 -d 5m cases/post-restaurant-grpc.js
