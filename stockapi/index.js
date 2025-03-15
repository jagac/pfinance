import express from 'express';
import yahooFinance from 'yahoo-finance2';

const app = express();
const port = 4000;

// Route to get stock price
app.get('/stock/:symbol', async (req, res) => {
    const symbol = req.params.symbol.toUpperCase();

    try {
        const quote = await yahooFinance.quote(symbol);
        const price = quote.regularMarketPrice;

        if (!price) {
            return res.status(404).json({ error: 'Stock price not found' });
        }

        res.json({
            symbol,
            price: price,
            timestamp: new Date(),
        });
    } catch (err) {
        res.status(500).json({ error: 'Failed to retrieve stock data', details: err.message });
    }
});

app.listen(port, () => {
    console.log(`Stock price service is running on http://localhost:${port}`);
});
