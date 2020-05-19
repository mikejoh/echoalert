#!/bin/bash

warning_alert='[
  {
    "labels": {
       "alertname": "DiskRunningFullSoon",
       "dev": "sda1",
       "instance": "instance01",
       "severity": "warning"
     },
     "annotations": {
        "info": "The disk sda1 is about to running full",
        "summary": "Please check sda1 and free up disk space"
      }
  }
]'

critical_alert='[
  {
    "labels": {
       "alertname": "DiskIsFull",
       "dev": "sda1",
       "instance": "instance01",
       "severity": "super-critical"
     },
     "annotations": {
        "info": "The disk sda1 is full",
        "summary": "Free up disk space!"
      }
  }
]'

curl -XPOST -d"$warning_alert" http://localhost:9093/api/v1/alerts
curl -XPOST -d"$critical_alert" http://localhost:9093/api/v1/alerts