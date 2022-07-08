const express = require('express');
const app = express();
const port = 3000;

const mapping = {
    '192.168.81.63': 'tv0',
    '192.168.93.19': 'tv4',
    '192.168.85.179': 'tv5'
}

app.get('/', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');

    const ipAddress = req.ip.startsWith('::ffff:') ? req.ip.substring(7) : req.ip;
    console.log('ip:', ipAddress);
    res.send({ ip: ipAddress, name: mapping[ipAddress] });
});

app.get('/time', (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    const d = new Date();
    // console.log('alo', d);
    res.send({ time: d });
});

app.listen(port, () => {
    console.log(`Example app listening on port ${port}`);
});