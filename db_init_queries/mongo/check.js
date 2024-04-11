const env = require('./db_init_queries/mongo/load_env.js');

db = db.getSiblingDB(env.MONGO_DB);

const result = db[env.MONGO_COLLECTION].find();
printjson(result);