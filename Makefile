windows-build:
	GOOS=windows GOARCH=amd64 go build -o miroir.exe

linux-build:
	GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o miroir
