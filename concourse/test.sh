set -e

project_path=$GOPATH/src/github.com/Bo0mer/
mkdir -p $project_path
cp -r dp $project_path

cd $project_path
go test -race ./...
