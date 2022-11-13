#!/bin/bash

set -o nounset   # disallow usage of unset vars  ( set -u )
set -o errexit   # Exit immediately if a pipeline returns non-zero.  ( set -e )
set -o errtrace  # Allow the above trap be inherited by all functions in the script.  ( set -E )
set -o pipefail  # Return value of a pipeline is the value of the last (rightmost) command to exit with a non-zero status
IFS=$'\n\t'      # Set $IFS to only newline and tab.

set -o functrace

cd "$(dirname "$0")/homebrew"

cp dops.rb dops_patch.rb


version="$(cd ../../../ && git tag --sort=-v:refname | grep -P 'v[0-9\.]' | head -1 | cut -c2-)"
cs0="$(cd ../../../ && sha256sum _out/dops_macos-amd64 | cut -d ' ' -f 1)"

echo "Version: ${version} (${cs0})"

sed --regexp-extended  -i "s/<<version>>/${version}/g"  dops_patch.rb
sed --regexp-extended  -i "s/<<shahash>>/${cs0}/g"      dops_patch.rb

cd ../../../
git clone https://github.com/Mikescher/homebrew-tap.git _out/homebrew-tap

cp "_data/package-data/homebrew/dops_patch.rb" _out/homebrew-tap/dops.rb
rm "_data/package-data/homebrew/dops_patch.rb"


cd _out/homebrew-tap/

git add dops.rb

if [ -z "$(git status --porcelain)" ]; then 
  echo "(!) Nothing changed -- nothing to commit"
else 
  git commit -m "dops v${version}"
fi



# git push manually (!)
