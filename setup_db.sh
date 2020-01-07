#!/bin/bash
echo Wait for servers to be up
sleep 10

HOSTPARAMS="--host cockroachdb --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

$SQL -e "CREATE USER IF NOT EXISTS scraping;"
$SQL -d scraping -e "CREATE DATABASE IF NOT EXISTS scraping WITH ENCODING = 'UTF8';"
$SQL -d scraping -e "GRANT ALL ON DATABASE scraping TO scraping;"