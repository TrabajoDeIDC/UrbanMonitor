# Urban Server

## Setup

This project requires no setup

## Usage

Run the program using:
```sh
go run UrbanServer.go [-test]
```
If you type `-test` you will run the program in TEST mode, meaning that random data will be created in a temp database for testing purposes

## Building

You can build the program using:
```sh
go build UrbanServer.go
```

and run it using:
```sh
./UrbanServer [-test]
```

## Production vs debug mode

Either by running the program using go or by building it and running it as a binary, you can set the production mode using `export GIN_MODE=release` and the debug mode using `export GIN_MODE=release`. If you don't specify the mode, the debug mode is set by default.