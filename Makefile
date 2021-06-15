
BINARY_PATH = ./dist/$(version)
BINARY_NAME = easy-mail

build: tidy clean darwin linux windows

tidy:
	go mod tidy
	find . -name "*.go" -type f -not -path "./vendor/*" | xargs -n1 go fmt

dist:
	mkdir -p ${BINARY_PATH}

darwin: dist
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ${BINARY_PATH}/${BINARY_NAME}
	cd ${BINARY_PATH} && zip ${BINARY_NAME}-darwin64.zip ${BINARY_NAME} && rm ${BINARY_NAME}

linux: dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${BINARY_PATH}/${BINARY_NAME}
	cd ${BINARY_PATH} && zip ${BINARY_NAME}-linux64.zip ${BINARY_NAME} && rm ${BINARY_NAME}

windows: dist
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${BINARY_PATH}/${BINARY_NAME}.exe
	cd ${BINARY_PATH} && zip ${BINARY_NAME}-win64.zip ${BINARY_NAME}.exe && rm ${BINARY_NAME}.exe

clean:
	rm -rf ${BINARY_PATH}

