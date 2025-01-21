@echo off
IF "%1"=="build" (
    go build -o bin/encryption-tool.exe src/main.go
    exit /b
)

IF "%1"=="run" (
    go run src/main.go %2 %3 %4 %5 %6
    exit /b
)

IF "%1"=="test" (
    go test ./...
    exit /b
)

IF "%1"=="clean" (
    if exist bin\encryption-tool.exe del /F bin\encryption-tool.exe
    exit /b
)

echo Usage:
echo run.bat build - Build the project
echo run.bat run - Run the project
echo run.bat test - Run tests
echo run.bat clean - Clean build files
