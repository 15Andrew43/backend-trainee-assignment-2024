#!/bin/bash


# load .env
source .env

export PGPASSWORD=$POSTGRES_PASSWORD

# postgres
echo "CREATING TABLES IN POSTGRES"
psql -h localhost -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB -a -f db_init_queries/postgres/create_tables.sql
sleep 1


# mongo
echo "CREATING COLLECTION IN MONGO"
mongosh --host localhost --port $MONGO_PORT --eval 'const script = load("db_init_queries/mongo/create_banner_data_collection.js")'
sleep 1
