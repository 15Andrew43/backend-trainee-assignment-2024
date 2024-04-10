db = db.getSiblingDB("mydatabase");

const result = db.banner_data.find();
printjson(result);