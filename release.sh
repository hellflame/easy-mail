#!/bin/bash

go build
export version=`./easy-mail -v`
make build
