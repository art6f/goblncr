import express from 'express';
import { statsMiddleware } from './logging/statsMiddleware';
import RouterRandom from './routes/randomDelay';

const app = express();

app.use(statsMiddleware)

app.use('/', RouterRandom);

const app_port = 8083;

app.listen(app_port, () => console.log("Express Test Server has started at port: " + app_port));

