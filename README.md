[![REUSE status](https://api.reuse.software/badge/github.com/SAP-samples/cf-autoscaler-custom-metrics-mtls)](https://api.reuse.software/info/github.com/SAP-samples/cf-autoscaler-custom-metrics-mtls)

# Demonstrate Scaling Metrics in Application Autoscaler

The SAP Business Technology Platform (BTP) provides a runtime environment for running your applications at scale.
This repository include a sample Golang application that use the Cloud Foundry runtime for SAP BTP.

Application Autoscaler, an autoscaling service in Cloud Foundry Runtime, automatically adds or removes application instances/nodes based on the workload. This can be done by defining an autoscaling policy, containing scaling rules.
This repository explains about how to send custom metrics to CF Autoscaler Service using mTLS authentication

### Application Autoscaler Use Cases
- Is your Cloud Foundry (CF) application unable to cope with large amounts of requests during peak hours?
- Have you seen a decrease in application performance in high-traffic conditions?
- Are you looking to reduce your runtime costs when there is less traffic?
- Do you want to scale your CF application based on your standard metrics (cpu, cpuutil, memory, throughput, reponsetime) or user-defined metrics ?



This is a sample repo containing golang application, demonstrating the following user cases:

- Send custom metrics to Cloud Foundry Autoscaler Service
- Use CPU Utilization to Scale Out and Scale In

## Practical Examples

- [Demonstrate Custom Metric MTLS Feature in Application Autoscaler](docs/custom-metrics-mtls.md)
- [Demonstrate CPU Utilization Metric Feature in Application Autoscaler](docs/cpu-utilization.md)


## License
Copyright (c) 2023 SAP SE or an SAP affiliate company. All rights reserved. This project is licensed under the Apache Software License, version 2.0 except as noted otherwise in the [LICENSE](LICENSE) file.
