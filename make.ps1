$command = $args[0]

switch ($command) {
    "clean" {
        Remove-Item .\build
    }
    "lint" {
        gofumpt -l -w .
	
        go vet
        go mod tidy
        go clean

        golangci-lint run
    }
    default { 
        New-Item -ItemType Directory -Path build -Force
        go build -o build
    }
}