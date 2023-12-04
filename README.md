# OIDC CloudWatch DataSource

## Introduction
This plugin is a data source for Grafana that uses OIDC (OpenID Connect) to access CloudWatch without AWS IAM user authentication credentials.  
If Grafana is running on GCP, it attempts to issue an STS id token using the Unique ID of the associated service account.

## Screenshots

![Configuration](https://github.com/arabian9ts/oidc-cloudwatch-datasource/raw/main/src/img/config.jpg)

![QueryEditor](https://github.com/arabian9ts/oidc-cloudwatch-datasource/raw/main/src/img/query.jpg)

## Support
- [x] CloudWatch Metric
- [ ] CloudWatch Logs
- [ ] CloudWatch Alarm
- [x] Grafana Alert
