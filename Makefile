build:
	go build -o bin/main main.go

run:
	go run main.go


compile:
	# Compiling for every OS Platform.

	# Windows
	# Windows x86 64 bit
	GOOS=windows GOARCH=amd64 go build -o bin/ffxiv-bard-windows-amd64.exe main.go
	# Windows ARM 64 bit
	GOOS=windows GOARCH=arm64 go build -o bin/ffxiv-bard-windows-arm64.exe main.go

	# MacOS
	# MacOS Darwin x86 64 bit
	GOOS=darwin GOARCH=amd64 go build -o bin/ffxiv-bard-darwin-amd64 main.go
	# MacOS Darwin ARM 64 bit
	GOOS=darwin GOARCH=arm64 go build -o bin/ffxiv-bard-darwin-arm64 main.go

	# Linux
	# Linux x86 64 bit
	GOOS=linux GOARCH=amd64 go build -o bin/ffxiv-bard-linux-amd64 main.go
	# Linux ARM 64 bit
	GOOS=linux GOARCH=arm64 go build -o bin/ffxiv-bard-linux-arm64 main.go
