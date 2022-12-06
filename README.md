# SpaceMoon 
An out-of-this-world online retail space 

## How to build
Run `go build -o spacemoon ./server` from the root directory to build the server executable (it will be created on the 
same root directory, with the name `server`). On windows, use `go build -o spacemoon.exe ./server` 
(it will create `spacemoon.exe`).

Alternatively, you can just run the server without building the executable by running `go run ./server`

Right now, it is using PORT 1234 to receive the requests.

---

# Changelog    

## [Unreleased]
### Added
* PayPal payment processing
* Stripe payment processing
* Instant Pay payment processing
* Google Pay payment processing
* Login and authentication token for private API calls

##  [v0.1.0] - _2022 12 03_
### Added
* Basic CRUD for category handling.

##  [v0.0.1] - _2022 12 02_
### Added
* Basic CRUD for product and main HTTP server.