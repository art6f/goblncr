import Router, { Request, Response } from 'express';
import { calculateDelay, randomCode } from '../common';

const DELAY_MIN = 10;
const DELAY_MAX = 400;

const code = 200;

const router = Router()

// generate random HTTP response code within the range 100-504
router.all('/', async (req: Request, res: Response) => {

    const delay = calculateDelay(DELAY_MIN, DELAY_MAX)
    await new Promise(resolve => setTimeout(resolve, delay));
    res.status(code).send(`OK, delay was: ${delay}ms, range: ${DELAY_MIN}-${DELAY_MAX}ms`);
});

// generate random code response with the defined delay
const delayRoutePt1 = '/delay/:delay1([1-9]+[\\d]{0,})';
router.all(delayRoutePt1, async (req: Request, res: Response) => {
    const delay = parseInt(req.params['delay1']);
    await new Promise(resolve => setTimeout(resolve, delay)); // sleep

    res.status(code).send(`Delay ${delay} ms, Response code ${code}`);
});

// generate random code response with the random delay within the defined range
router.all(`${delayRoutePt1}-:delay2([1-9]+[\\d]+)`, async (req: Request, res: Response) => {
    let delayFrom = parseInt(req.params['delay1']);
    let delayTo = parseInt(req.params['delay2']);
    let delay = calculateDelay(delayFrom, delayTo);

    await new Promise(resolve => setTimeout(resolve, delay));
    res.status(code).send(`Delay ${delay}ms, Response code ${code}, range: ${delayFrom}-${delayTo}ms`);
});

export default router;