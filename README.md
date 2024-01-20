# check-value-from-metrics-plugin

## Description

Transform Mackerel metrics plugin to check plugin.

## Synopsis

```sh
mackerel-plugin-anyone | check-value-from-metrics-plugin -target anyone.value -gt 100
```

## Installation

```sh
sudo mkr plugin install Arthur1/check-value-from-metrics-plugin
```

## Setting for mackerel-agent

```
[plugin.checks.check-value-from-mackerel-plugin-anyone]
command = "mackerel-plugin-anyone | /opt/mackerel-agent/plugins/bin/check-value-from-metrics-plugin -target anyone.value -gt 100"
```

## Usage

### Options

```
  -critical
        create WARNING alert when target value meets condition (default true)
  -gt float
        set condition: [target value] > [option value] (default NaN)
  -lt float
        set condition: [target value] < [option value] (default NaN)
  -name string
        checker name for report (default "check-value-from-metrics-plugin")
  -strict
        create UNKNOWN alert when target metric is not found
  -target string
        target metric key
  -warning
        create CRITICAL alert when target value meets condition
```
