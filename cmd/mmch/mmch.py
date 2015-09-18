import operator
import sys

defline = sys.stdin.readline()
s = sys.stdin.readline().strip()

count = {}
for b in s:
   count[b] = 1 + count.get(b, 0)
nA = count['A']
nC = count['C']
nG = count['G']
nU = count['U']

def matches(m, n):
   return reduce(operator.mul, range(n - m + 1, n + 1), 1)

mAU = matches(min(nA, nU), max(nA, nU))
mGC = matches(min(nC, nG), max(nC, nG))
m = mAU * mGC
print m
