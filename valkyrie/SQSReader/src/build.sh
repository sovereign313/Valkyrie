export GOPATH=$(pwd | sed 's/src//')
rm -rf ./github.com
rm -rf ./golang.com
rm -rf ./gopkg.in
go get
GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-w' sqsreader.go 
if [[ $? -ne 0 ]]; then
        echo "Failed To Compile"
        exit
fi

