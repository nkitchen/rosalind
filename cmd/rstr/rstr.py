import sys

line = sys.stdin.readline()
words = line.split()
N = int(words[0])
pGC = float(words[1])
pAT = 1 - pGC

line = sys.stdin.readline()
s = line.strip()

count = {}
for b in s:
   count[b] = 1 + count.get(b, 0)
kA = count['A']
kC = count['C']
kG = count['G']
kT = count['T']

k = len(s)

ps = 0.5**k * pGC**(kC + kG) * pAT**(kA + kT)
ps_once = 1 - (1 - ps)**N
print "{:.3}".format(ps_once)
