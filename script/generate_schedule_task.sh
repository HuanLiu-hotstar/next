#!/bin/bash
# gdate is date cmd in MacOS
if [ $# -lt 1 ]; then
    echo "./$0 tab-split-file" # can export as tab split file in google sheet
    exit 1
fi
delimeter="|"
while read line; do

    if [ -z "$line" ]; then
        continue
    fi
    line=$(echo "$line" | awk -F "\t" '{printf("%s %s|%s %s|%s\n",$5,$6,$5,$7,$9)}')
    #echo "line:"$line
    start_date=$(echo "$line" | awk -F $delimeter '{print $1}')
    end_date=$(echo "$line" | awk -F $delimeter '{print $2}')
    name=$(echo $line | awk -F $delimeter '{print $3}')
    #echo "-----$start_date,  $end_date, $name"
    start_time=$(gdate -d "$start_date 3 hour ago" '+%Y-%m-%dT%H:%M:%S')
    end_time=$(gdate -d"$end_date 3 hour" '+%Y-%m-%dT%H:%M:%S')
    #echo "$start_time, $end_time"
    echo "- {"
    echo "      name: \"$name\","
    echo "      level: \"CRICKET\","
    echo "      start_time: \"$start_time+05:30\","
    echo "      end_time: \"$end_time+05:30\","
    echo "  }"
done <$1
