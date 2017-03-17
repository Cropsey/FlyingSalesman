#!/usr/bin/env python

from __future__ import print_function
import os
import subprocess
import csv

try:
    out = subprocess.check_output("go build".split())
    out += subprocess.check_output("go build fspcmd/main.go".split())
except:
    print("Couldn't build app!")
    print(out)
    os.exit(1)


data_sets = [
        '/tmp/data_20.txt',
        '/tmp/data_50.txt',
        '/tmp/data_100.txt',
        ]

params = {
        'DCFS_MAX_BRANCHES' : range(1, 10, 2),
        'DCFS_DISC_W'       : [x * 0.01 for x in range(0, 200, 25)],
        'DCFS_MIN_DISC'     : [x * 0.01 for x in range(-50, 50, 10)],
        'DCFS_NEXT_AVG_W'   : [x * 0.01 for x in range(0, 100, 10)], 
        'DCFS_DISC_THRESH'  : range(0, 1000, 200),
        }

results = {}
for d in data_sets:
    results[d] = {}
    for p in params.keys():
        results[d][p] = []

os.environ["FSP_ENGINE"] = "DCFS"
for p, vals in params.items():
    for d in data_sets:
        for v in vals:
            os.environ[p] = str(v)
            print("Testing {} with {} = {}...".format(d, p, v), end='')
            out = subprocess.check_output("./main -v".split(), stdin=open(d), stderr=open("/dev/null"))
            del os.environ[p]
            price = int(out.split("\n")[0])
            print(price)
            results[d][p].append(price)

del os.environ["FSP_ENGINE"]

print(results)
