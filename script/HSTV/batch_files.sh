#!/bin/bash

current_date=`date +%Y-%m-%d`

for((i=0;i<7;i++))
do
    d=`date -v+"$i"d +%m-%d-%Y`
    python3 generate_schedule_with_duration.py "${d}"
    echo "Running batch file for HSTV $i"
done