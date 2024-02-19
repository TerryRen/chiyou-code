@echo off
cd /d %~dp0
echo "==================  clear ==================="
for %%f in (bin\*) do if not "%%~xf"==".yml" del "%%f"
echo "=====================mmc====================="
echo "================= mmc windows ==============="
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o bin/mmc_amd64.exe   main.go
echo "================= mmc darwin  ==============="
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin/mmc_darwin_amd64 main.go
echo "================= mmc linux   ==============="
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/mmc_linux_amd64  main.go
