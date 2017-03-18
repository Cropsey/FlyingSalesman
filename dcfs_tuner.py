#!/usr/bin/env python

from __future__ import print_function
import os
import subprocess
import csv
import Queue
import threading

thread_pool_size = 4

def worker(q, results):
    while True:
        params = q.get()
        if params is None:
            return
        d, p, v,  = params
        env = dict(os.environ)
        env[p] = str(v)
        print("Testing {} with {} = {}...".format(d, p, v))
        out = subprocess.check_output("./main -v".split(), stdin=open(d), stderr=open("/dev/null"), env=env)
        #del os.environ[p]
        price = int(out.split("\n")[0])
        print("{}:{}={}\t{}".format(d, p, v, price))
        results[d][p].append((v, price))


def build_app():
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


def main():
    results = {}
    for d in data_sets:
        results[d] = {}
        for p in params.keys():
            results[d][p] = []
    
    q = Queue.Queue()
    threads = [threading.Thread(target=worker, args=(q, results, )) for i in range(thread_pool_size)]
    os.environ["FSP_ENGINE"] = "DCFS"
    for p, vals in params.items():
        for d in data_sets:
            for v in vals:
                q.put((d, p, v, ))
    
    for t in threads:
        t.start()
        q.put(None)

    for t in threads:
        t.join()
    
    del os.environ["FSP_ENGINE"]
    
    print(results)
    
if __name__ == "__main__":
    build_app()
    main()

