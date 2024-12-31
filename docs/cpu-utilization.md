# Demonstrate CPU Utilization Metric Feature in Application Autoscaler

Cloud Foundry (CF) is a powerful platform that allows developers to deploy and scale applications with ease. To further enhance the scalability of applications, a new dynamic scaling metric called CPU Utilization (cpuutil) has been added to Application Autoscaler. This metric is linked to the CPU entitlement usage of an app, making it easier for developers to scale their applications based on CPU usage.

## Introducing CPU Utilization (cpuutil) metric in Application Autoscaler

The newly introduced CPU utilization metric (cpuutil) in Application Autoscaler allows users to define scaling rules based on CPU usage in percentage. With cpuutil, developers no longer need to manually calculate the CPU entitlement of their app and adjust the thresholds accordingly. Instead, they can simply make use of cpuutil metric and set thresholds from 0% to 100%.

How CPU Utilization Metric is different from CPU metric

Previously, scaling applications using CPU metric require calculating the CPU entitlement manually. This requires manual adjustment and careful monitoring of CPU usage thresholds. Also, the calculated thresholds (defined in the scaling policy) do not guarantee the expected results as the CPU entitlement can be increased or decreased (good vs bad apps - discussed later) by the platform. This manual adjustment process was time-consuming and often prone to human error.
The CPU Utilization aka cpuutil metric makes user’s life easier by defining CPU utilization as a percentage ranging from 0 to 100. This also eliminates the need of calculating CPU entitlement explicitly by the app developer.

In summary, CPU usage shows the actual consumption, while CPU entitlement indicates the allocated limit based on application’s memory.

Checkout the [blog](https://community.sap.com/t5/technology-blogs-by-sap/application-scalability-with-cpu-utilization-metric-in-cloud-foundry/ba-p/13696223) for more details on how to use CPU Utilization metric in Application Autoscaler.