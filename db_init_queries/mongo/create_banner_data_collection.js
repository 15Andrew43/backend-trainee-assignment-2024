const env = require('./db_init_queries/mongo/load_env.js');

db = db.getSiblingDB(env.MONGO_DB);

db.createCollection(env.MONGO_COLLECTION);
