#!/bin/sh

set -e
cd "$( cd "$(dirname "$0")"; pwd -P)"

git submodule update
/usr/local/bin/hugo
git add .
git commit -m 'Update website'
git push

git checkout master
rm -rIf *
git checkout dev -- public
mv public/* .
rmdir public
git add .
git commit -m 'Update website'
git push

git checkout -
