#!/bin/bash

# $1 - number of monitors
# $2 - "stretch" or "fit"
# $3 - filename

res="$(($1 * 1080)):1920"
ext="${3##*.}"
bname=`basename "$3" ".$ext"`

if [ "$3" = fit ]
then
    extrarg=":force_original_aspect_ratio=decrease,pad=$res:-1:-1:color=black"
else
    extrarg=""
fi

ffmpeg -y -i raw/"$3" -c:a copy -vf "scale=$res$extrarg,setsar=1:1" processed/"$bname.$ext" || exit

for i in `seq "$1"`
do
	ffmpeg -y -i processed/"$bname.$ext" -vf "crop=1080:1920:$((($i - 1)*1080)):0,transpose=dir=2" sdata/"$bname.$i.$ext" || exit
done
