#!/bin/bash

resolution=`ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 "$1"`

if ! [[ "$resolution" = 6480x1920 ]]
then
	echo Wrong resolution want 6480x1920, have "$resolution"
	exit
fi

for i in `seq 6`
do
	ffmpeg -i "$1" -filter:v crop=1080:1920:$((($i - 1)*1080)):0 `basename "$1" .mp4`.$i.mp4
	#echo crop=1080:1920:$((($i - 1)*1080)):0
done
