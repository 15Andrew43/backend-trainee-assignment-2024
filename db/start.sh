#!/bin/bash


# load .env
source .env

# postgres
echo "CREATING TABLES IN POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db/postgres/create_tables.sql
echo "INSERTING TEST DATA IN POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db/postgres/inserts.sql
echo
echo

# mongo
echo "CREATING COLLECTION IN MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db/mongo/create_banner_data_collection.js")'
echo "INSERTING TEST DATA IN MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db/mongo/inserts.js");'
echo
echo


echo "CHECKING POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db/postgres/check.sql
echo
echo

echo "CHECKIING MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db/mongo/check.js");'
echo
echo