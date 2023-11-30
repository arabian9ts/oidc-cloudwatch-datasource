import React, { useEffect, useState } from 'react';
import { Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { EditorField, EditorFieldGroup, EditorRow, EditorRows } from '@grafana/experimental';
import { DataSource } from '../datasource';
import { DataSourceOptions, Query, Namespace, MetricName } from '../types';
import { Dimensions } from './Dimensions';

type Props = QueryEditorProps<DataSource, Query, DataSourceOptions>;

const statistics = ['Average', 'Sum', 'Maximum', 'Minimum', 'SampleCount', 'IQM'];

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const [namespaces, setNamespaces] = useState<Namespace[]>([]);
  const [metricNames, setMetricNames] = useState<MetricName[]>([]);

  useEffect(() => {
    datasource.listNamespaces().then(setNamespaces).catch(console.error);
  }, [datasource, setNamespaces]);

  useEffect(() => {
    if (query.namespace === undefined || query.namespace === '') {
      return;
    }
    datasource.listMeticNames(query.namespace).then(setMetricNames).catch(console.error);
  }, [datasource, setMetricNames, query.namespace]);


  const onNamespaceChange = (value: string) => {
    onChange({ ...query, namespace: value });
    onRunQuery();
  };

  const onMetricNameChange = (value: string) => {
    onChange({ ...query, metricName: value });
    onRunQuery();
  };

  const onStatisticChange = (value: string) => {
    onChange({ ...query, statistic: value });
    onRunQuery();
  };

  const listNamespacesAsOptions = (): Array<SelectableValue<string>> => {
    return [
      ...namespaces.map((o) => {
        return { value: o.name, label: o.name };
      }),
    ];
  };

  const listMetricNamesAsOptions = (): Array<SelectableValue<string>> => {
    return [
      ...metricNames.map((o) => {
        return { value: o.name, label: o.name };
      }),
    ];
  };

  const statisticsAsOptions = (): Array<SelectableValue<string>> => {
    return [
      ...statistics.map((o) => {
        return { value: o, label: o };
      }),
    ];
  };

  return (
    <>
      <EditorRows>
        <EditorRow>
          <EditorFieldGroup>
            <EditorField label="Namespace" width={26}>
              <Select
                options={listNamespacesAsOptions()}
                value={query.namespace}
                width={28}
                onChange={(e) => onNamespaceChange(e?.value!)}
                className="inline-element"
                isClearable={true}
              />
            </EditorField>
            <EditorField label="MetricName" width={16}>
              <Select
                options={listMetricNamesAsOptions()}
                value={query.metricName}
                width={28}
                onChange={(e) => onMetricNameChange(e?.value!)}
                className="inline-element"
                isClearable={true}
              />
            </EditorField>

            <EditorField label="Statistic" width={16}>
              <Select
                allowCustomValue
                value={query.statistic}
                options={statisticsAsOptions()}
                onChange={(e) => onStatisticChange(e?.value!)}
              />
            </EditorField>
          </EditorFieldGroup>
        </EditorRow>

        <EditorRow>
          <EditorField label="Dimensions">
            <Dimensions
              onChange={onChange}
              onRunQuery={onRunQuery}
              datasource={datasource}
              query={query}
            />
          </EditorField>
        </EditorRow>
      </EditorRows>
    </>
  );
}
