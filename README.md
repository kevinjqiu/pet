```
            _
 _ __   ___| |_
| '_ \ / _ \ __|
| |_) |  __/ |_
| .__/ \___|\__|
|_|
```

PET - Prometheus Expression Testing framework

A utility for unit testing Prometheus Query and Rules

Usage
-----

    pet [testfiles]

e.g.,

    pet test*
    pet tests/*.yaml

Writing Test Cases
------------------

Take a look at the sample testcase [here](https://github.com/kevinjqiu/pet/tree/master/tests/testcase.yaml).

Attribution
-----------

This tool is inspired by [query-tester](https://github.com/m-lab/prometheus-support/tree/master/cmd/query_tester) from m-lab
