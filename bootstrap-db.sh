#!/bin/bash
docker exec -it stubber-db psql -U postgres -d postgres -f /docker-entrypoint-initdb.d/dump.sql
