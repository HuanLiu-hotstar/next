contentId=$1

duration=$(
    
    curl --location "https://origin-multi-get.preprod.hotstar-labs.com/o/v2/multi/get/content?ids=$contentId" \
    --header 'x-platform-code: PCTV' \
    --header 'x-country-code: in' \
    --header 'accept: */*' \
    --header 'sec-fetch-mode: cors' | jq ".body.results.map.\"$contentId\".duration"
)
echo "$contentId,$duration"

