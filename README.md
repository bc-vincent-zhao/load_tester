# load_tester

This is an simple aggregator/automation tool to facilitate load testing on API endpoints.

Its main features include:
- automatically generate random dummy urls and request bodies with different payload size to better simulate real life traffic and measure longest path rather than constantly hitting a cached path somewhere
- `./load_test saturate` invoke `wrk` to roughly estimate the max throughput of the endpoint
- `./load_test histogram` invoke `wrk2` with constant request rate to determine the response time in high precision (99.9999%)
- `./load_test ts` invoke `vegeta` with constatn request rate to generate latency data in time series raw format, which to be used to verify the results of the above 2 and allow custom analysis of the raw data 

To see all the available options, simply run
```
$ load_tester
```
