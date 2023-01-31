# MOODY SKELTON

## Getting Started
1. Make sure you have [Go](https://go.dev) installed.
2. Clone the repo
```bash
git clone 
```
3. Go to the directory and run go mod tidy to add missing requirements and to drop unused requirements
```bash
cd moody && go mod tidy
```
3. Setup your .env file
```bash
cp .env-example .env && vi .env
```
4. Start
```bash
go run main.go
```
## Build for production
1. Compile packages and dependencies
```bash
go build -o moody main.go
```
2. Setup .env file for production
```bash
cp .env-example .env && vi .env
```
3. Run executable file with systemd, supervisor, pm2 or other process manager
```bash
./moody
```
