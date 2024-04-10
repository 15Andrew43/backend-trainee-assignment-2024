db = db.getSiblingDB("mydatabase");

db.createCollection("banner_data", {
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
