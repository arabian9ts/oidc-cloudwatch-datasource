import React, { useEffect, useState } from 'react';
import { Select } from '@grafana/ui';
import { QueryEditorProps, toOption } from '@grafana/data';
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
    if (!query.namespace) {
      return;
    }
    datasource.listMeticNames(query.namespace).then(setMetricNames).catch(console.error);
  }, [datasource, setMetricNames, query.namespace]);

  const onQueryChange = <Key extends keyof Query, Value extends Query[Key]>(
    option: Key,
    value: Value | undefined,
  ) => {
    onChange({ ...query, [option]: value });
    onRunQuery();
  };

  return (
    <>
      <EditorRows>
        <EditorRow>
          <EditorFieldGroup>
            <EditorField label="Namespace" width={26}>
              <Select
                options={namespaces.map((n) => toOption(n.name))}
                value={query.namespace}
                width={28}
                onChange={(e) => onQueryChange('namespace', e?.value!)}
                className="inline-element"
                isClearable={true}
              />
            </EditorField>
            <EditorField label="MetricName" width={16}>
              <Select
                options={metricNames.map((m) => toOption(m.name))}
                value={query.metricName}
                width={28}
                onChange={(e) => onQueryChange('metricName', e?.value!)}
                className="inline-element"
                isClearable={true}
              />
            </EditorField>

            <EditorField label="Statistic" width={16}>
              <Select
                allowCustomValue
                value={query.statistic}
                options={statistics.map((s) => toOption(s))}
                onChange={(e) => onQueryChange('statistic', e?.value!)}
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
