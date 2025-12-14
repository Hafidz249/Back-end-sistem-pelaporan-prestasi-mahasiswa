#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."

# Install swag if not exists
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate docs
swag init -g main.go -o ./docs --parseDependency --parseInternal

echo "Swagger documentation generated successfully!"
echo "Access documentation at: http://localhost:8080/swagger/index.html"