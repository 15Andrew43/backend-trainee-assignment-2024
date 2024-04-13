#!/bin/bash


# load .env
source .env

# postgres
export PGPASSWORD=$POSTGRES_PASSWORD
echo "CREATING TABLES IN POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db_init_queries/postgres/create_tables.sql
sleep 1
echo "INSERTING TEST DATA IN POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db_init_queries/postgres/inserts.sql
echo
echo

# mongo
echo "CREATING COLLECTION IN MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db_init_queries/mongo/create_banner_data_collection.js")'
sleep 1
echo "INSERTING TEST DATA IN MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db_init_queries/mongo/inserts.js");'
echo
echo

sleep 1
echo "CHECKING POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db_init_queries/postgres/check.sql
echo
echo

echo "CHECKIING MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db_init_queries/mongo/check.js");'
echo
echo