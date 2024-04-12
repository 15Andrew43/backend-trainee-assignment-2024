const env = require('./db_init_queries/mongo/load_env.js');

db = db.getSiblingDB(env.MONGO_DB);

db[env.MONGO_COLLECTION].insertMany([
    { id: 101, content: "Banner Data 1" },
    { id: 102, content: "Banner Data 2" },
    { id: 103, content: "Banner Data 3" },
    { id: 104, content: {"title": "some_title", "text": "some_text", "url": "some_url"} },
]);