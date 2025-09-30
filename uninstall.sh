#!/bin/env bash

CHK="/usr/local/bin/run" && [ ! -f "$CHK" ] && echo "Not installed" && exit
DEL="/usr/local/bin/run" && sudo rm -f $DEL && echo "Uninstalled"
