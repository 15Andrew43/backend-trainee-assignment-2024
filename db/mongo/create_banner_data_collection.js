const env = require('./db/mongo/load_env.js');

db = db.getSiblingDB(env.MONGO_DB);

db.createCollection(env.MONGO_COLLECTION, {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["id", "content"],
            properties: {
                id: {
                    bsonType: "int",
                    description: "must be an integer and is required"
                },
                content: {
                    bsonType: "string",
                    description: "must be a string and is required"
                }
            }
        }
    }
});
