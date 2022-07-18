#!/bin/bash

# $1 - number of monitors
# $2 - resolution
# $3 - "stretch" or "fit"
# $4 - filename

ext="${4##*.}"
bname=`basename "$4" ".$ext"`

if [ "$3" = fit ]
then
    extrarg=":force_original_aspect_ratio=decrease,pad=$2:-1:-1:color=black"
else
    extrarg=""
fi

ffmpeg -i raw/"$4" -c:a copy -vf "scale=$2$extrarg,setsar=1:1" processed/"$bname.$1.$3.$ext" || exit

for i in `seq "$1"`
do
	ffmpeg -i processed/"$bname.$1.$3.$ext" -vf "crop=1080:1920:$((($i - 1)*1080)):0,transpose=dir=2" sdata/"$bname.$1.$3.$i.$ext" || exit
done
