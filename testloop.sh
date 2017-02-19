#!/bin/bash

export PATH=$PATH:$GOPATH/bin

tmux new-session \; set status off \; split-window -hp 20 redgreen \; last-pane
