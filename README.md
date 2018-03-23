# s0_exporter

Prometheus exporter for a s0 reader.
Converts impulses to Watt:

watt := (countedImpulses / wattPerPulse) * 1000 * (secondsBetweenScrapes)

```
# HELP s0_power Power used between scrapes (in watt)
# TYPE s0_power gauge
s0_power 39
# HELP s0_time Time since last scrape (in seconds)
# TYPE s0_time gauge
s0_time 5.32
```
