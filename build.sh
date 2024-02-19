echo "=====================mmc====================="
echo "==================== clear ===================="
rm -rf bin/mmc_*
echo "================= mmc windows ==============="
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/mmc_amd64.exe    main.go
echo "================= mmc darwin  ==============="
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o bin/mmc_darwin_amd64 main.go
echo "================= mmc linux   ==============="
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o bin/mmc_linux_amd64  main.go
