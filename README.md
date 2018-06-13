```
                                 _       _            _
 _ __  _ __ ___  _ __ ___   __ _| |     | |_ ___  ___| |_
| '_ \| '__/ _ \| '_ ` _ \ / _` | |_____| __/ _ \/ __| __|
| |_) | | | (_) | | | | | | (_| | |_____| ||  __/\__ \ |_
| .__/|_|  \___/|_| |_| |_|\__, |_|      \__\___||___/\__|
|_|                           |_|

```

Utility for unit testing Prometheus Query and Rules

Usage
-----

    promql-test [testfiles]

e.g.,

    promql-test test*
    promql-test tests/*


Test File
---------

Test file uses the Prometheus internal test file DSL.

For a tutorial, see [here](https://github.com/m-lab/prometheus-support/blob/master/cmd/query_tester/README.md)

For sample test files, see [here](https://github.com/prometheus/prometheus/blob/master/promql/testdata/)

Attribution
-----------

This tool is inspired by [query-tester](https://github.com/m-lab/prometheus-support/tree/master/cmd/query_tester) from m-lab
