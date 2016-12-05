set -e

cd dp
version=$(cat .version)
echo -e "v$version" > release_name
echo -e "v$version" > release_tag

git log `git describe --tags --abbrev=0`..HEAD --oneline > release_body
