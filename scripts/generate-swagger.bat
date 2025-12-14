@echo off
echo Generating Swagger documentation...

REM Check if swag is installed
swag version >nul 2>&1
if %errorlevel% neq 0 (
    echo Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate docs
swag init -g main.go -o ./docs --parseDependency --parseInternal

echo Swagger documentation generated successfully!
echo Access documentation at: http://localhost:8080/swagger/index.html
pause