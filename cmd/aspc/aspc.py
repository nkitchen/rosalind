#!/usr/bin/env python

import sys

n = int(sys.argv[1])
m = int(sys.argv[2])

def binom(n, k):
   b = 1
   for x in xrange(1, k + 1):
      b *= n - x + 1
      assert b % x == 0
      b //= x
   return b

print sum(binom(n, k) for k in range(m, n + 1))
