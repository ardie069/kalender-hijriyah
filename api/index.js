const express = require("express");
const app = express();

app.get("/api/hello", (req, res) => {
    res.json({ message: "Hello from Vercel API!" });
});

module.exports = app;
