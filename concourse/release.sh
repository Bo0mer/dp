set -ex

mkdir -p build/
build_dir=$(cd build && pwd)

cd dp
version=$(cat .version)
echo -e "v$version" > $build_dir/release_name
echo -e "v$version" > $build_dir/release_tag

git log `git describe --tags --abbrev=0`..HEAD --oneline > $build_dir/release_body

cd ..

project_path=$GOPATH/src/github.com/Bo0mer/
mkdir -p $project_path
cp -r dp $project_path

cd $project_path/dp

GOOS=linux GOARCH=amd64 go build -v -o $build_dir/dp_linux_amd64
GOOS=darwin GOARCH=amd64 go build -v -o $build_dir/dp_darwin_amd64
GOOS=windows GOARCH=amd64 go build -v -o $build_dir/dp_windows_amd64.exe
