#!/bin/bash

# contentID=$1
# jq -n --arg content_id $contentID '{"contentId":$content_id,"token":"","country":"in","platform":"android","skipCddsQuery":false,"clientCapabilities":""}'

if [ $# -lt 1 ]; then
    echo "bash $0 content_id_list(1540039986 1540040496)"
    exit -1
fi

for contentID in $*
do
    uri=`bash fetchContentDetail.sh $contentID 2>/dev/null | jq -r '.data.fetchContentDetail.contentMeta.coreAttributes.images.horizontalImage.url'`
    echo "HSTV Poster, $contentID, https://img1.hotstarext.com/image/upload/$uri"
done
