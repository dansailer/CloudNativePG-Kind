#!/bin/bash
# File: start-multi-view.sh
tmux new-session -d -s multi-view
tmux split-window -h
tmux select-pane -t 0
tmux split-window -v
tmux select-pane -t 2
tmux split-window -v
tmux send-keys -t 0 "k9s" C-m
tmux send-keys -t 1 "watch -n 5 kubectl cnpg status" C-m
#tmux send-keys -t 2 "./start-processing.sh" C-m
tmux attach-session -t multi-view