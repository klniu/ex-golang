#!/bin/bash
#/d/mplus/git_x86/bin/find.exe /e/home/mplus/data/mon/trunk/server/src/tagbar -type f -mmin 0
#/d/mplus/git_x86/bin/find.exe . -type f -print0 | xargs -0 stat
watch() {
#echo watching folder $1/ every $2 secs.
touch -d "`date +%c`" watch.txt

while [[ true ]]
do
	timestamp="`date +%c`"
    files=`/d/mplus/git_x86/bin/find.exe $1 -type f -cnewer watch.txt`
    #files=`/d/mplus/git_x86/bin/find.exe $1 -type f -mmin 0`
	touch -d "$timestamp" watch.txt
    if [[ $files == "" ]] ; then
        echo "nothing changed"
    else
        echo "changed", $files
    fi
    sleep 1
done
}

watch "/e/home/mplus/data/mon/trunk/server/src/tagbar" 0