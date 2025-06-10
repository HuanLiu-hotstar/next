#!/bin/bash

# contentID=$1
# jq -n --arg content_id $contentID '{"contentId":$content_id,"token":"","country":"in","platform":"android","skipCddsQuery":false,"clientCapabilities":""}'

if [ $# -lt 1 ]; then
    echo "bash $0 content_id_list(1540039986 1540040496)"
    exit -1
fi

# i=1
# for contentID in $*
# do
#     data=`bash fetchContentDetail.sh $contentID 2>/dev/null `
#     uri=` echo $data | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.images.horizontalImage.url'`
#     title=`echo $data| jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.title' `
#     description=`echo $data | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.description'`
#     echo "$i. HSTV, $title, $contentID, https://img1.hotstarext.com/image/upload/$uri"

#     jq -n --arg content_id $contentID --arg uri "$uri" --arg channelName "$title" --arg channelDescription "$description" '{"id":$content_id,"ChannelImageURL":$uri,"channelNumber":"channelNumber","channelDescription":$channelDescription,"channelName":$channelName}'
#     i=$((i+1))
# done


function get_json() {
    contentID=$1
    channelID=$2
    data=`bash fetchContentDetail.sh $contentID 2>/dev/null `
    uri=` echo $data | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.images.horizontalImage.url'`
    title=`echo $data| jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.title' `
    description=`echo $data | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.description'`
    # echo "$i. HSTV, $title, $contentID, https://img1.hotstarext.com/image/upload/$uri"
    
    res=`jq -n --arg content_id $contentID --arg uri "$uri" --arg channelName "$title" --arg channelDescription "$description" --arg channelID $channelID '{"id":$content_id,"ChannelImageURL":$uri,"channelNumber":$channelID,"channelDescription":$channelDescription,"channelName":$channelName}'`
    echo $res
}


i=0
json_array="[]"
while read line
do
    if [ $i -eq 0 ]; then
        i=$((i+1))
        continue
    fi
    
    contentID=`echo $line | awk -F, '{print $1}'`
    channelID=`echo $line | awk -F, '{print $2}'`
    echo $contentID, $channelID
    out=`get_json $contentID $channelID`
    echo $out
    json_array=`echo $json_array | jq --argjson obj "$out" '. + [$obj]'`
    i=$((i+1))
done < $1

echo $json_array | jq .