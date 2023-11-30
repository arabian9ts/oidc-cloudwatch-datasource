import { css, cx } from '@emotion/css';
import React, { useState, useEffect } from 'react';
import { QueryEditorProps, toOption, SelectableValue } from '@grafana/data';
import { Select, useTheme2 } from '@grafana/ui';
import { AccessoryButton, Space } from '@grafana/experimental';
import { DataSource } from '../datasource';
import { DataSourceOptions, Query, Dimension } from '../types';

type Props = QueryEditorProps<DataSource, Query, DataSourceOptions>;

export const Dimensions = ({ query, datasource, onChange, onRunQuery }: Props) => {
  const [dimensions, setDimensions] = useState<Dimension[]>([]);

  useEffect(() => {
    if (!query.namespace || !query.metricName) {
      return;
    }
    datasource.listDimensions(query.namespace, query.metricName).then(setDimensions).catch(console.error);
  }, [datasource, setDimensions, query.namespace, query.metricName]);

  const theme = useTheme2();
  const add = (name: string, value: string) => {
    query.dimensions.push({ name, value });
    onChange(query);
    onRunQuery();
  }
  const remove = (index: number) => {
    query.dimensions.splice(index, 1);
    onChange(query);
    onRunQuery();
  }
  const update = (index: number, name: string, value: string) => {
    query.dimensions[index] = { name, value };
    onChange(query);
    onRunQuery();
  }

  return (
    <>
      {query.dimensions.map((item, index) => (
        <>
          <Select
            width="auto"
            value={item.name ? toOption(item.name) : null}
            showAllSelectedWhenOpen={true}
            allowCustomValue
            options={uniqOptions(dimensions.map((o) => o.name))}
            onChange={(e) => { update(index, e.value || '', item.value) }}
          />
          <span
            className={cx(css({
              padding: theme.spacing(0, 1),
              alignSelf: 'center',
            }))}
          >
            =
          </span>
          <Select
            width="auto"
            value={item.value ? toOption(item.value) : null}
            showAllSelectedWhenOpen={true}
            allowCustomValue
            options={uniqOptions(dimensions.filter((o) => o.name === item.name).map((o) => o.value))}
            onChange={(e) => { update(index, item.name, e.value || '') }}
          />
          <AccessoryButton icon="times" variant="secondary" onClick={() => remove(index)} type="button" />
          <Space h={3} v={2} layout='inline' />
        </>
      ))}

      <Space v={2} layout='block' />
      <AccessoryButton icon="plus" variant="secondary" onClick={() => { add('', '') }} type="button">
        Add dimension
      </AccessoryButton>
    </>
  );
};

const uniqOptions = (opts: string[]): Array<SelectableValue<string>> => {
  return [...new Set(opts)].map((v) => toOption(v));
}
