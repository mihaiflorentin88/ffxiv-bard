build:
	go build -o bin/main main.go

run:
	go run main.go


compile:
	# Compiling for every OS Platform.

	# Windows
	# Windows x86 64 bit
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/ffxiv-bard-windows-amd64.exe main.go

	# Windows ARM 64 bit
	#CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CGO_CFLAGS="-target aarch64-w64-mingw32" go build -o bin/ffxiv-bard-windows-arm64.exe main.go

	# MacOS
	# MacOS Darwin x86 64 bit
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/ffxiv-bard-darwin-amd64 main.go
	# MacOS Darwin ARM 64 bit
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o bin/ffxiv-bard-darwin-arm64 main.go

	# Linux
	# Linux x86 64 bit
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/ffxiv-bard-linux-amd64 main.go
	# Linux ARM 64 bit
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o bin/ffxiv-bard-linux-arm64 main.go
