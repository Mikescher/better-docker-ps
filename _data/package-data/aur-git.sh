#!/bin/bash

set -o nounset   # disallow usage of unset vars  ( set -u )
set -o errexit   # Exit immediately if a pipeline returns non-zero.  ( set -e )
set -o errtrace  # Allow the above trap be inherited by all functions in the script.  ( set -E )
set -o pipefail  # Return value of a pipeline is the value of the last (rightmost) command to exit with a non-zero status
IFS=$'\n\t'      # Set $IFS to only newline and tab.

cd "$(dirname "$0")/aur-git"
git clean -ffdX

version=$(cd ../../../ && git tag --sort=-v:refname | grep -P 'v[0-9\.]' | head -1 | cut -c2-)

echo "Version: ${version}"

sed --regexp-extended  -i "s/pkgver=[0-9\.]+/pkgver=${version}/g" PKGBUILD



namcap PKGBUILD
makepkg --printsrcinfo > .SRCINFO
makepkg


cd ../../../
pwd
git clone ssh://aur@aur.archlinux.org/dops-git.git _out/dops-git
cp _data/package-data/aur-git/PKGBUILD _out/dops-git/PKGBUILD
cp _data/package-data/aur-git/.SRCINFO _out/dops-git/.SRCINFO


cd _out/dops-git

git add PKGBUILD
git add .SRCINFO

if [ -z "$(git status --porcelain)" ]; then 
  echo "(!) Nothing changed -- nothing to commit"
else 
  git commit -m "v${version}"
fi


cd "../../_data/package-data/aur-git"
git clean -ffdX

# git push manually (!)
