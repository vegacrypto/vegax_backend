#!/bin/sh

HOME="/home/env/go"
ENV="prod"

export DALINK_GO_CONFIG_PATH="$HOME/conf/param-$ENV.yaml"
export DALINK_TPL_PATH="$/home/tpl/"

nohup ./router-linux &
