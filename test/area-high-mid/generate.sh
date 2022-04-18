#!/bin/bash
set -e
if [ $# -lt 1 ]; then
    echo "usage:./$0 area_filename"
fi

cat $1 | jq '.data.middlelist' | jq '.[].area_name' -r | awk '{print $1,$2}' | uniq >tmpmid.txt
cat $1 | jq '.data.highlist' | jq '.[].area_name' -r | awk '{print $1,$2}' | uniq >tmphigh.txt

#city-code, match city with code
citycode=city-code.txt

function getdata() {

    while read line; do
        time=$(date +%Y-%m-%d)
        city=$(echo $line | awk '{print $2}')
        prov=$(echo $line | awk '{print $1}')
        count=$(echo "$prov" | grep "市" | wc -l)
        printcity=$city
        if [ $count -gt 0 ]; then
            city=$prov
            printcity=""
        fi

        code=$(grep $city $citycode | awk '{print $1}')
        level="$2"
        echo "$code,$level,$prov,$printcity,$time"
    done <$1

}

echo "qhdm,risk,sfmc,dsmc,gxsj"
getdata tmphigh.txt "高"
getdata tmpmid.txt "中"
