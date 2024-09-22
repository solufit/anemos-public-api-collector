#! /bin/bash

# Run the migration

cd /go/src/app

echo "Createing DataBase Scheme"
dbmate --url $DATABASE_URL -s db/scheme.sql up

echo "Migrating DataBase"
dbmate --url $DATABASE_URL -d db/seed -s db/scheme.sql up
