#!/usr/bin/env bash

rm tarolas-cov
rm report.out
go test -coverpkg="./..." -c -tags testrunmain -o tarolas-cov
./tarolas-cov -test.run "^TestRunMain$" -test.coverprofile=report.out
go tool cover -html=report.out -o coverage.html
go tool cover -func=report.out > a.txt
