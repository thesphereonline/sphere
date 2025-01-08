const express = require('express')
const bodyParser = require('body-parser')
const cors = require('cors')

const app = express();
const PORT = 3001;

app.use(bodyParser.json());
app.use(cors());

app.get('/', (req, res) => {
    res.json({ message: "Welcome to The Sphere Backend!" });
})

app.listen(PORT, () => {
    console.log(`Backend running on http://localhost:${PORT}`);
})