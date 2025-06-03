#!/bin/bash

# contentID=$1
# jq -n --arg content_id $contentID '{"contentId":$content_id,"token":"","country":"in","platform":"android","skipCddsQuery":false,"clientCapabilities":""}'

if [ $# -lt 1 ]; then
    echo "bash $0 content_id_list(1540039986 1540040496)"
    exit -1
fi

i=1
for contentID in $*
do
    data=`bash fetchContentDetail.sh $contentID 2>/dev/null `
    uri=` echo $data | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.images.horizontalImage.url'`
    title=`echo $data| jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.title' `
    echo "$i. HSTV, $title, $contentID, https://img1.hotstarext.com/image/upload/$uri"
    i=$((i+1))
done
