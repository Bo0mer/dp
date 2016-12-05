set -ex

mkdir -p build/
build_dir=$(cd build && pwd)

version=$(cat version/number)
echo -e "v$version" > $build_dir/release_name
echo -e "$version" > $build_dir/release_tag

cd dp
git log `git describe --tags --abbrev=0`..HEAD --oneline > $build_dir/release_body
cd ..


project_path=$GOPATH/src/github.com/Bo0mer/
mkdir -p $project_path
cp -r dp $project_path

cd $project_path/dp

GOOS=linux GOARCH=amd64 go build -v -o $build_dir/dp_linux_amd64
GOOS=darwin GOARCH=amd64 go build -v -o $build_dir/dp_darwin_amd64
GOOS=windows GOARCH=amd64 go build -v -o $build_dir/dp_windows_amd64.exe
