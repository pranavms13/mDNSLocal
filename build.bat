go build -o build/mDNSLocal-win-amd64.exe mDNSLocal

set GOOS=linux
set GOARCH=amd64
go build -o build/mDNSLocal-linux-amd64 mDNSLocal

set GOOS=darwin
set GOARCH=amd64
go build -o build/mDNSLocal-darwin-amd64 mDNSLocal