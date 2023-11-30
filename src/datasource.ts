import { DataSourceInstanceSettings, CoreApp } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { Query, DataSourceOptions, DEFAULT_QUERY, Namespace, MetricName, Dimension } from './types';

export class DataSource extends DataSourceWithBackend<Query, DataSourceOptions> {
  constructor(
    instanceSettings: DataSourceInstanceSettings<DataSourceOptions>,
  ) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<Query> {
    return DEFAULT_QUERY;
  }

  listNamespaces(): Promise<Namespace[]> {
    return this.getResource<Namespace[]>('api/namespaces');
  }

  listMeticNames(namespace: string): Promise<MetricName[]> {
    return this.getResource<MetricName[]>(`api/metricNames?namespace=${namespace}`);
  }

  listDimensions(namespace: string, metricName: string): Promise<Dimension[]> {
    return this.getResource<Dimension[]>(`api/dimensions?namespace=${namespace}&metricName=${metricName}`);
  }
}
