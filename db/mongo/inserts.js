db = db.getSiblingDB("mydatabase");

db.banner_data.insertMany([
    { id: 101, content: "Banner Data 1" },
    { id: 102, content: "Banner Data 2" },
    { id: 103, content: "Banner Data 3" }
]);