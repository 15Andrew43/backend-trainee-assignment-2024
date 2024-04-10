#!/bin/bash

# postgres
echo "CREATING TABLES IN POSTGRES"
psql -h localhost -p 5432 -U myuser -d mydb -a -f db/postgres/create_tables.sql
echo "INSERTING TEST DATA IN POSTGRES"
psql -h localhost -p 5432 -U myuser -d mydb -a -f db/postgres/inserts.sql
echo
echo

# mongo
echo "CREATING COLLECTION IN MONGO"
mongosh --host localhost --port 27017 --eval 'const script = load("db/mongo/create_banner_data_collection.js")'
echo "INSERTING TEST DATA IN MONGO"
mongosh --host localhost --port 27017 --eval 'const script = load("db/mongo/inserts.js");'
echo
echo


echo "CHECKING POSTGRES"
psql -h localhost -p 5432 -U myuser -d mydb -a -f db/postgres/check.sql
echo
echo

echo "CHECKIING MONGO"
mongosh --host localhost --port 27017 --eval 'const script = load("db/mongo/check.js");'
echo
echo