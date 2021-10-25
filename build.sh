#! /bin/bash

# ember s apiHost=https://demo1.dev:5001
# go run edition/community.go -port=5001 -forcesslport=5002 -cert selfcert/cert.pem -key selfcert/key.pem -salt=tsu3Acndky8cdTNx3

NOW=$(date)
source ~/.profile
echo "Build process started $NOW"

echo "Building Svelte assets..."
cd frontend
npm run build
cd ..

echo "Copying Svelte assets..."
rm -rf assets
mv frontend/dist/assets ./
mv frontend/dist/index.html templates

echo "Compiling for Linux..."
env GOOS=linux GOARCH=amd64 go build
#echo "Compiling for macOS Intel..."
#env GOOS=darwin GOARCH=amd64 go build -mod=vendor -trimpath -o bin/documize-community-darwin-amd64 ./edition/community.go
#echo "Compiling for macOS ARM..."
#env GOOS=darwin GOARCH=arm64 go build -mod=vendor -trimpath -o bin/documize-community-darwin-arm64 ./edition/community.go
#echo "Compiling for Windows..."
#env GOOS=windows GOARCH=amd64 go build -mod=vendor -trimpath -o bin/documize-community-windows-amd64.exe ./edition/community.go
#echo "Compiling for ARM..."
#env GOOS=linux GOARCH=arm go build -mod=vendor -trimpath -o bin/documize-community-linux-arm ./edition/community.go
#echo "Compiling for ARM64..."
#env GOOS=linux GOARCH=arm64 go build -mod=vendor -trimpath -o bin/documize-community-linux-arm64 ./edition/community.go

echo "Finished."

# CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo
# go build -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo test.go
# ldd test
