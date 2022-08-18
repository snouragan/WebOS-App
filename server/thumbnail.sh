#!/bin/bash

ext="${1##*.}"
bname=`basename "$1" ".$ext"`

ffmpeg -y -ss 00:00:01.00 -i "raw/$1" -vf 'scale=320:320:force_original_aspect_ratio=decrease' -vframes 1 "data/$1.thumb.jpg"
