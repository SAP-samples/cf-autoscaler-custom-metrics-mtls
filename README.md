# Golang App - Demonstrate Custom Metrics Feature in Application Autoscaler

This is a sample Golang app to showcase

- How to send custom metrics to [CF Autoscaler Service](https://github.com/cloudfoundry/app-autoscaler-release) using mTLS authentication

## Build App
### Option 1

`env GOOS=linux GOARCH=amd64 go build -o build/golang-autoscaler-custom-metrics`

### Option 2
`make build`

## Run Test

`make test`

## Deploy App

`cf push --no-start -f deploy/app-manifest.yml -p deploy/build`

```shell
$ cf push --no-start -f deploy/app-manifest.yml -p deploy/build
Pushing app golang-autoscaler-custom-metrics to org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...
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
    stack: cflinuxfs3
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
Creating service instance ak-test-autoscaler in org ak_autoscaler_dev / space ak-test-space as as ak-user@xxx.com...

Service instance ak-test-autoscaler created.
OK
```

## Bind App with service
`$ cf bind-service golang-autoscaler-custom-metrics ak-test-autoscaler -c deploy/custom-metrics-policy.json`

```shell
$ cf bind-service golang-autoscaler-custom-metrics ak-test-autoscaler -c deploy/custom-metrics-policy.json
Binding service instance ak-test-autoscaler to app golang-autoscaler-custom-metrics in org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...
OK

TIP: Use 'cf restage golang-autoscaler-custom-metrics' to ensure your env variable changes take effect

```

## check VCAP_SERVICE Env

```shell
$ cf env golang-autoscaler-custom-metrics
Getting env variables for app golang-autoscaler-custom-metrics in org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...
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

## Start an app
`cf start golang-autoscaler-custom-metrics`

```shell
cf start golang-autoscaler-custom-metrics
Starting app golang-autoscaler-custom-metrics in org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...

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

Starting app golang-autoscaler-custom-metrics in org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...

Waiting for app to start...

Instances starting...
Instances starting...
Instances starting...

name:              golang-autoscaler-custom-metrics
requested state:   started
routes:            golang-autoscaler-custom-metrics.cfapps.sap.hana.ondemand.com
last uploaded:     Fri 30 Dec 16:04:06 CET 2022
stack:             cflinuxfs3
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

## Send custom metrics to autoscaler

Scale out > `https://golang-autoscaler-custom-metrics.cfapps.abc.com/busy/301`

```shell
Every 2.0s: cf app golang-autoscaler-custom-metrics                                                                                                                                          XNQNV6VGJC: Fri Dec 30 16:19:22 2022

Showing health and status for app golang-autoscaler-custom-metrics in org ak_autoscaler_dev / space ak-test-space as ak-user@xxx.com...

name:              golang-autoscaler-custom-metrics
requested state:   started
routes:            golang-autoscaler-custom-metrics.cfapps.sap.hana.ondemand.com
last uploaded:     Fri 30 Dec 16:16:35 CET 2022
stack:             cflinuxfs3
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


## References:

- TLS Read
  - Recommend - https://youngkin.github.io/post/gohttpsclientserver/
  - https://smallstep.com/hello-mtls/doc/combined/nodejs/requests
  - https://smallstep.com/hello-mtls/doc/client/go

- Application development
  - https://levelup.gitconnected.com/a-practical-approach-to-structuring-go-applications-7f77d7f9c189
  - GIN WebFramework -  https://blog.logrocket.com/building-microservices-go-gin/
  - https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin
  - https://blog.logrocket.com/building-microservices-go-gin/

