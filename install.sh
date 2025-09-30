#!/bin/env bash

CHK="/usr/local/bin/run" && [ -f "$CHK" ] && echo "Updating..." && sudo rm -f $CHK
RUN="/usr/local/bin" && sudo cp ./run $RUN && echo "Done"
