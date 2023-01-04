# Changelog

## [Unreleased]
### Fixed
* 1 username per person

---

##  [v1.2.4] - _2023 01 04_
### Fixed
* enable CORS on POST for direct message handler

##  [v1.2.3] - _2023 01 04_
### Fixed
* send DMs only to existing users (check first)

##  [v1.2.2] - _2023 01 04_
### Fixed
* clean "Bearer" from Bearer token if needed

##  [v1.2.1] - _2023 01 03_
### Fixed
* small fixes on profile handler

##  [v1.2.0] - _2023 01 03_
### Added
* Direct Messages

---

##  [v1.1.0] - _2022 12 28_
### Added
* user profile
    * GET profile
    * PUT profile
    * add profile if not exists (and user does exist)
    * user profile persistence

---

##  [v1.0.4] - _2022 12 27_
### Fixed
* you can only erase your own posts

##  [v1.0.3] - _2022 12 27_
### Added
* delete posts

##  [v1.0.2] - _2022 12 23_
### Fixed
* likes saved

##  [v1.0.1] - _2022 12 23_
### Added
* check if user exists before creating it
### Fixed
* changed post uri to string for Google storage
* added error handling on file saving

##  [v1.0.0] - _2022 12 22_
### Added
* Save image and ad URLs to post

---

##  [v0.6.3] - _2022 12 18_
### Fixed
* Fixed CORS on social media Posts API

##  [v0.6.2] - _2022 12 18_
### Added
* TLS added to API

##  [v0.6.1] - _2022 12 18_
### Added
* Google Cloud Persistence for Social Network Posts

##  [v0.6.0] - _2022 12 18_
### Added
* Basic Social Network Posts

---

##  [v0.5.1] - _2022 12 15_
### Fixed
* Login endpoint was not checking for credentials

##  [v0.5.0] - _2022 12 14_
### Added
* Product seller based on the user adding the product

---

##  [v0.4.0] - _2022 12 14_
### Added
* Login and Product Persistence

---

##  [v0.3.0] - _2022 12 09_
### Added
* Basic ratings on `product/rating`

---

##  [v0.2.3] - _2022 12 07_
### Fixed
* CORS pre-flight response on login handler
* Added allowed methods for product and category to CORS pre-flight response

##  [v0.2.2] - _2022 12 07_
### Fixed
* Enabled CORS origin header on all requests methods response headers

##  [v0.2.1] - _2022 12 07_
### Added
* Enabled CORS pre-flight response

##  [v0.2.0] - _2022 12 06_
### Added
* Login and authentication token for private API calls
* All methods except GET protected for products
* All methods except GET protected for categories
* Updated README instructions

---

##  [v0.1.0] - _2022 12 03_
### Added
* Basic CRUD for category handling.

---

##  [v0.0.1] - _2022 12 02_
### Added
* Basic CRUD for product and main HTTP server.