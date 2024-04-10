const env = require('./db/mongo/load_env.js');

db = db.getSiblingDB(env.MONGO_DB);

const result = db[env.MONGO_COLLECTION].find();
printjson(result);