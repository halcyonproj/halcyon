#!/bin/bash

# run halcyon itself
make run &

CONTROLLERS=$@
pushd ..

for controller in ${CONTROLLERS}; do
    pushd ${controller}
    make run &
    popd
done
wait
