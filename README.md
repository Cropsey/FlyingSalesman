# Flying Salesman Problem

[![Build Status](https://travis-ci.org/Cropsey/fsp.svg?branch=master)](https://travis-ci.org/Cropsey/fsp)

Kiwi challenge - https://travellingsalesman.cz/

## Arguments

* `-v` be verbose and output a lot of stuff to stderr
* `-s` just read input and print some statistics about the data
* `-i int` set timeout for the solution (default 30s)

## Env vars

* `FSP_ENGINE` selects engine to solve the problem, possible values are: `DCFS`, `SITM`, `MITM`, `RANDOM`, `BHDFS`, `BN`, `GREEDY`, `ROUNDS`
* `DCFS_MAX_BRANCHES` branching limit for DCFS engine
* `DCFS_DISC_W` discount contribution factor to flight evaluation
* `DCFS_NEXT_AVG_W` next node avg flight price contribution to flight evaluation
* `DCFS_MIN_DISC` minimal discount needed to consider flight (`0.3` means 30% discount, `-0.2` means 20% overpriced flight)
* `DCFS_DISC_THRESH` minimal price to apply minimal discount rule

