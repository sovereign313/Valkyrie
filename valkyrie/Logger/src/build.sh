export GOPATH=$(pwd | sed 's/src//')
rm -rf ./github.com
go get
GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-w' logger.go 
if [[ $? -ne 0 ]]; then
        echo "Failed To Compile"
        exit
fi