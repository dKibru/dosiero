#! /bin/bash

NOW=$(date)
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
GIN_MODE=release
env GOOS=linux GOARCH=amd64 go build
mv dosiero releases/dosiero-linux

echo "Compiling for Mac ..."
GIN_MODE=release
env GOOS=darwin GOARCH=arm64 go build
mv dosiero releases/dosiero-mac

echo "Finished."
