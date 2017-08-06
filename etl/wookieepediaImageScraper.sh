#!/usr/bin/env bash

wookieepediaUrl=$1
contains=${2:-http}

# get all the links found on page with 'vignette' and '.jpeg' or '.png'
# apply optional filter
# strip beginning up until 'http'
# strip end starting with '/revision'
# output to json like {images:[]}

lynx -dump -hiddenlinks=listonly \
${wookieepediaUrl} \
| grep 'vignette.*\.jp\|vignette.*\.png' \
| grep ${contains} \
| sed -n 's/^.*http/http/p' \
| sed 's/\/revision.*//' \
| awk ' BEGIN { ORS = ""; print "["; } { print "\/\@"$0"\/\@"; } END { print "]"; }' \
| sed "s^\"^\\\\\"^g;s^\/\@\/\@^\", \"^g;s^\/\@^\"^g" \
| jq '{images: .}'
