#!/usr/bin/env python

import operator
import sys

n = int(sys.argv[1])
k = int(sys.argv[2])

r = range(n - k + 1, n + 1)
print reduce(operator.mul, r)
