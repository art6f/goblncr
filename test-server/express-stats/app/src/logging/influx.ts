import { InfluxDB, Point, WriteApi } from '@influxdata/influxdb-client';

const influx = new InfluxDB({
  url: process.env.INFLUX_URL!,
  token: process.env.INFLUX_TOKEN!,
});

const writeApi: WriteApi = influx.getWriteApi(
  process.env.INFLUX_ORG!,
  process.env.INFLUX_BUCKET!
);

interface Stat {
  podId: string,
  timestamp: Date;
  endpoint: string;
  statusCode: number;
  responseTimeMs: number;
}

export default function writeLog(stat: Stat): void {
  console.log(`Respose from [${stat.podId}]: endpoint ${stat.endpoint}`)

  const point = new Point('request')
    .tag('pod_id', stat.podId)
    .tag('endpoint', stat.endpoint)
    .intField('status_code', stat.statusCode)
    .intField('response_time', stat.responseTimeMs)
    .timestamp(stat.timestamp);

  writeApi.writePoint(point);
}
