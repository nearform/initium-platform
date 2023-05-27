#!/bin/bash

if [ $# -ne 3 ]
  then
    echo "Usage: ./upload-assets.sh <GITHUB_REPOSITORY> <FILE_NAME> <GITHUB_REF>"
    exit 1
fi

ACTION_REPOSITORY=$1
ACTION_FILE=$2
ACTION_TAG=$(echo "$3" | sed -e s/refs\\/tags\\///g | sed -e s/refs\\/heads\\///g)
RELEASE_RESPONSE=$(curl \
-H "Accept: application/vnd.github+json" \
-H "Authorization: Bearer $GITHUB_TOKEN" \
https://api.github.com/repos/"$ACTION_REPOSITORY"/releases/tags/"$ACTION_TAG")

ASSET_LIST=$(curl \
-H "Accept: application/vnd.github+json" \
-H "Authorization: Bearer $GITHUB_TOKEN" \
https://api.github.com/repos/"$ACTION_REPOSITORY"/releases/$(echo "$RELEASE_RESPONSE" | jq .id)/assets)

echo "Pre-existing asset list: $(echo "$ASSET_LIST" | jq)"
i=0
for f in $(echo "$ASSET_LIST" | jq -r '.[]? | .name'); do
if [ "$f" = "$ACTION_FILE" ] ; then
    echo "This file already exists, we'll overwrite it"
    curl \
    -X DELETE \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    https://api.github.com/repos/"$ACTION_REPOSITORY"/releases/assets/$(echo "$ASSET_LIST" | jq -r --argjson i "$i" '.[$i].id')
fi
(( i+= 1 ))
done

curl \
-X POST \
-H "Accept: application/vnd.github+json" \
-H "Authorization: Bearer $GITHUB_TOKEN" \
-H "Content-Type: application/yaml" \
-H "Content-Length: $(wc -c "$ACTION_FILE" | awk '{print $1}')" \
--data-binary "@$ACTION_FILE" \
"$(echo "$RELEASE_RESPONSE" | jq .assets_url | sed -e s/api.github.com/uploads.github.com/g | sed -e s/\"//g)?name=$ACTION_FILE&label=$ACTION_FILE";

echo "File $ACTION_FILE uploaded to release $ACTION_TAG successfully"
exit 0
