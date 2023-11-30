import React from 'react';
import { InlineField, Input, Select } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { DataSourceOptions, TokenIssuer, TokenIssuerType } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> { }

export function ConfigEditor(props: Props) {
  const labelWidth = 25;
  const valueWidth = 60;
  const { onOptionsChange, options } = props;
  const { jsonData } = options;

  const onOptionChange = <Key extends keyof DataSourceOptions, Value extends DataSourceOptions[Key]>(
    option: Key,
    value: Value | undefined,
  ) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, [option]: value },
    });
  };

  const issuerAsOptions = (): Array<SelectableValue<TokenIssuer>> => {
    return Object.entries(TokenIssuerType).map(([_, value]) => {
      return { value: value, label: value };
    });
  };

  return (
    <>
      <InlineField
        label="Access Token Issuer"
        labelWidth={labelWidth}
        tooltip="Access Token Issuer such as GoogleServiceAccount"
        required
      >
        <Select
          options={issuerAsOptions()}
          onChange={(e) => onOptionChange('issuer', e?.value!)}
          value={jsonData.issuer}
          placeholder=""
          width={valueWidth}
        />
      </InlineField>
      <InlineField
        label="Credentials Path"
        labelWidth={labelWidth}
        tooltip="Path to credentials file. If not set, will use default credentials."
      >
        <Input
          onChange={(e) => onOptionChange('credentialsPath', e.currentTarget.value)}
          value={jsonData.credentialsPath}
          placeholder="/path/to/credentials.json"
          width={valueWidth}
        />
      </InlineField>
      <InlineField
        label="Assume Role ARN"
        labelWidth={labelWidth}
        tooltip="ARN of the role to assume."
        required
      >
        <Input
          onChange={(e) => onOptionChange('assumeRole', e.currentTarget.value)}
          value={jsonData.assumeRole}
          placeholder="arn:aws:iam:*"
          width={valueWidth}
        />
      </InlineField>
      <InlineField
        label="STS Region"
        labelWidth={labelWidth}
        tooltip="Region to use for STS."
        required
      >
        <Input
          onChange={(e) => onOptionChange('stsRegion', e.currentTarget.value)}
          value={jsonData.stsRegion}
          placeholder=""
          width={valueWidth}
        />
      </InlineField>
      <InlineField
        label="Monitoring Region"
        labelWidth={labelWidth}
        tooltip="Region to use for CloudWatch."
        required
      >
        <Input
          onChange={(e) => onOptionChange('region', e.currentTarget.value)}
          value={jsonData.region}
          placeholder=""
          width={valueWidth}
        />
      </InlineField>
    </>
  );
}
