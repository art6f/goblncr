import { Request, Response, NextFunction } from 'express';
import writeLog from './influx';

export function statsMiddleware(req: Request, res: Response, next: NextFunction): void {
  const start = Date.now();

  const podId = process.env.K8S_POD_ID || 'UNK';

  res.on('finish', () => {
    const stat = {
      podId: podId,
      endpoint: req.path,
      statusCode: res.statusCode,
      responseTimeMs: Date.now() - start,
      timestamp: new Date(),
    };

    writeLog(stat);
  });

  next();
}
