# Bubblegum 
A next-gen social media network and app 

## How to build
Made with Go.

Run `go build -o spacemoon ./server` from the root directory to build the server executable (it will be created on the 
same root directory, with the name `server`). On windows, use `go build -o spacemoon.exe ./server` 
(it will create `spacemoon.exe`).

Alternatively, you can just run the server without building the executable by running `go run ./server`

Right now, it is using PORT 1234 to receive the requests.
