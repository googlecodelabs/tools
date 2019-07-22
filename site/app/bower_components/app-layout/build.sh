#!/bin/bash
#

#git checkout master
#git pull
#git checkout gh-pages
#git reset --hard master
#git push -f

a=`find . -name "docs.html" -or -name "index.html" -not -path "*/test/*"`
b=`find . -name "*.html" ! -name 'sample-content.html' -path "*/demo/*"`

c=(`for R in "${a[@]}" "${b[@]}" ; do echo "$R" ; done | sort -du`)

for f in ${c[@]}; do
  echo "vulcanize " $f
  vulcanize --inline-css --inline-scripts $f > $f.build
  mv $f.build $f
done

#git commit -a -m "build"
#git push
