SHELL=cmd

# build: builds all binaries
build: clean build_server
	@echo All binaries built!

# clean: cleans all binaries and runs go clean
clean:
	@echo Cleaning...
	@echo y | DEL /S dist
	@go clean
	@echo Cleaned and deleted binaries

# build_server: builds the server
build_server:
	@echo Building server...
	@go build -o dist/gostripe.exe ./server
	@echo Server built!

# start: start the server
start: start_server

# start_server: starts the server
start_server: build_server
	@echo Starting the server...
	${SHELL} /c start /B ./dist/gostripe.exe
	@echo Server running!

# stop: stop the server
stop: stop_server
	@echo Server stopped

# stop_server: stops the server
stop_server: 
	@echo Stopping the server...
	@taskkill /IM gostripe.exe /F
	@echo Server stopping