import { DataQuery, DataSourceJsonData } from '@grafana/schema';

export interface Query extends DataQuery {
  namespace: string;
  metricName: string;
  statistic: string;
  dimensions: Dimension[];
}

export interface Dimension {
  name: string;
  value: string;
}

export type Namespace = {
  name: string;
};

export type MetricName = {
  name: string;
  value: string;
};

export const DEFAULT_QUERY: Partial<Query> = {
  namespace: '',
  metricName: '',
  statistic: '',
  dimensions: [] as Dimension[],
};

export const TokenIssuerType = {
  Google: 'GOOGLE',
} as const;

export type TokenIssuer = typeof TokenIssuerType[keyof typeof TokenIssuerType];

export interface DataSourceOptions extends DataSourceJsonData {
  issuer: TokenIssuer;
  credentialsPath?: string;
  assumeRole: string;
  stsRegion: string;
  region: string;
}
