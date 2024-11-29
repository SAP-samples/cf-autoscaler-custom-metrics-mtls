[![REUSE status](https://api.reuse.software/badge/github.com/SAP-samples/cf-autoscaler-custom-metrics-mtls)](https://api.reuse.software/info/github.com/SAP-samples/cf-autoscaler-custom-metrics-mtls)


# Demonstrate Scaling Metrics in Application Autoscaler

This sample golang application demonstrate the following user cases:

- Send custom metrics to Cloud Foundry Autoscaler Service
- Use CPU Utilization to Scale Out and Scale In

## Description

The SAP Business Technology Platform (BTP) provides a runtime environment for running your applications at scale.
This repository include a sample Golang application that use the Cloud Foundry runtime for SAP BTP.

Check out the step-by-step guide on [SAP Blogs](https://blogs.sap.com/?p=1613870&preview=true&preview_id=1613870)

Application Autoscaler, an autoscaling service in Cloud Foundry Runtime, automatically adds or removes application instances/nodes based on the workload. This can be done by defining an autoscaling policy, containing scaling rules.
This repository explains about how to send custom metrics to CF Autoscaler Service using mTLS authentication

### Application Autoscaler Use Cases
- Is your Cloud Foundry (CF) application unable to cope with large amounts of requests during peak hours?
- Have you seen a decrease in application performance in high-traffic conditions?
- Are you looking to reduce your runtime costs when there is less traffic?
- Do you want to scale your CF application based on your custom metrics (other than CPU, memory, throughput, reponsetime) ?

This tutorial includes the following steps
 - Build and Deploy sample Golang Application
 - Create Autoscaler Service Instance
 - Bind App with Autoscaler Service
 - Start Golang App
 - Send Custom Metrics to Autoscaler Service
 - Monitor Scale-Out and Scale-In 

## Requirements

- You should have access to Cloud Foundry Runtime

## Build and Deploy sample Golang Application

### Build App
#### Option 1 `env GOOS=linux GOARCH=amd64 go build -o build/golang-autoscaler-custom-metrics`
#### Option 2 `make build`

### Run Test

`make test`

### Deploy App

`cf push --no-start -f deploy/app-manifest.yml -p deploy/build`

sample-output
```shell
$ cf push --no-start -f deploy/app-manifest.yml -p deploy/build
Pushing app golang-autoscaler-custom-metrics to org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...
Applying manifest file deploy/app-manifest.yml...

Updating with these attributes...
  ---
  applications:
  - name: golang-autoscaler-custom-metrics
    disk-quota: 128M
    instances: 1
    path: /Users/development/golang-autoscaler-custom-metrics/deploy/build
    memory: 128M
+   default-route: true
    stack: cflinuxfs4
    buildpacks:
      binary_buildpack
    command: ./golang-autoscaler-custom-metrics
Manifest applied
All files found in remote cache; nothing to upload.
Waiting for API to complete processing files...

name:              golang-autoscaler-custom-metrics
requested state:   stopped
routes:            golang-autoscaler-custom-metrics.cfapps.sap.hana.ondemand.com
last uploaded:
stack:
buildpacks:

type:            web
sidecars:
instances:       0/1
memory usage:    128M
start command:   ./golang-autoscaler-custom-metrics
     state   since                  cpu    memory   disk     details
#0   down    2022-12-30T14:22:58Z   0.0%   0 of 0   0 of 0

```

## Create Autoscaler Service Instance

```shell
$ cf create-service autoscaler standard ak-test-autoscaler
Creating service instance ak-test-autoscaler in org ak_autoscalerxxx / space ak-test-space as as ak-user@xxx.com...

Service instance ak-test-autoscaler created.
OK
```
## Bind App with Autoscaler Service
`$ cf bind-service golang-autoscaler-custom-metrics ak-test-autoscaler -c deploy/custom-metrics-policy.json`

```shell
$ cf bind-service golang-autoscaler-custom-metrics ak-test-autoscaler -c deploy/custom-metrics-policy.json
Binding service instance ak-test-autoscaler to app golang-autoscaler-custom-metrics in org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...
OK

TIP: Use 'cf restage golang-autoscaler-custom-metrics' to ensure your env variable changes take effect

```
### Check VCAP_SERVICE Environment Variable

```shell
$ cf env golang-autoscaler-custom-metrics
Getting env variables for app golang-autoscaler-custom-metrics in org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...
System-Provided:
VCAP_SERVICES: {
  "autoscaler": [
    {
      "binding_guid": "c02fba90-83b2-4753-a31f-1d0827978632",
      "binding_name": null,
      "credentials": {
        "custom_metrics": {
          "mtls_url": "https://autoscaler-metrics-mtls.cf.xxx.com",
          "password": "xxxxxxxxx",
          "url": "https://autoscaler-metrics.cf.xxx.com",
          "username": "xxxxxxxxx"
        }
      },
      "instance_guid": "3eeac1a2-ced1-4748-a3b4-a87188a7d3c6",
      "instance_name": "ak-test-autoscaler",
      "label": "autoscaler",
      "name": "ak-test-autoscaler",
      "plan": "standard",
      "provider": null,
      "syslog_drain_url": null,
      "tags": [
        "autoscaler",
        "app-autoscaler",
        "cf-autoscaler"
      ],
      "volume_mounts": []
    }
  ]
}


VCAP_APPLICATION: {
  ...
  ....
  ...

```
## Start Golang App
`cf start golang-autoscaler-custom-metrics`


```shell
cf start golang-autoscaler-custom-metrics
Starting app golang-autoscaler-custom-metrics in org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...

Staging app and tracing logs...
   Downloading binary_buildpack...
   Downloaded binary_buildpack
   Cell 038ca894-a748-49d5-85eb-be08569dfe48 creating container for instance 95b8140e-065c-4116-8c42-085f3f33db27
   Cell 038ca894-a748-49d5-85eb-be08569dfe48 successfully created container for instance 95b8140e-065c-4116-8c42-085f3f33db27
   Downloading app package...
   Downloading build artifacts cache...
   Downloaded build artifacts cache (215B)
   Downloaded app package (6M)
   -----> Binary Buildpack version 1.0.47
   Exit status 0
   Uploading droplet, build artifacts cache...
   Uploading droplet...
   Uploading build artifacts cache...
   Uploaded build artifacts cache (216B)

Starting app golang-autoscaler-custom-metrics in org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...

Waiting for app to start...

Instances starting...
Instances starting...
Instances starting...

name:              golang-autoscaler-custom-metrics
requested state:   started
routes:            golang-autoscaler-custom-metrics.cfapps.sap.hana.ondemand.com
last uploaded:     Fri 30 Dec 16:04:06 CET 2022
stack:             cflinuxfs4
buildpacks:
	name               version   detect output   buildpack name
	binary_buildpack   1.0.47    binary          binary

type:           web
sidecars:
instances:      1/1
memory usage:   128M
     state     since                  cpu    memory      disk        details
#0   running   2022-12-30T15:04:16Z   0.0%   0 of 128M   0 of 128M

```

## Send Custom Metrics to Autoscaler Service

Scale out > `$ curl https://golang-autoscaler-custom-metrics.cfapps.abc.com/busy/301`

```shell
Every 2.0s: cf app golang-autoscaler-custom-metrics                                                                                                                                          XNQNV6VGJC: Fri Dec 30 16:19:22 2022

Showing health and status for app golang-autoscaler-custom-metrics in org ak_autoscalerxxx / space ak-test-space as ak-user@xxx.com...

name:              golang-autoscaler-custom-metrics
requested state:   started
routes:            golang-autoscaler-custom-metrics.cfapps.sap.hana.ondemand.com
last uploaded:     Fri 30 Dec 16:16:35 CET 2022
stack:             cflinuxfs4
buildpacks:
        name               version   detect output   buildpack name
        binary_buildpack   1.0.47    binary          binary

type:           web
sidecars:
instances:      2/2
memory usage:   128M
     state     since                  cpu    memory          disk            details
#0   running   2022-12-30T15:16:44Z   0.9%   23.6M of 128M   10.9M of 128M
#1   running   2022-12-30T15:18:13Z   0.9%   22.9M of 128M   10.9M of 128M

```

Scale In > `$ curl https://golang-autoscaler-custom-metrics.cfapps.abc.com/not-busy/190`

## License
Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This project is licensed under the Apache Software License, version 2.0 except as noted otherwise in the [LICENSE](LICENSE) file.
