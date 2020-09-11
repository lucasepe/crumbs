# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2020-11-09
### Added
- 🎉 new flag `-images-type` to specify a default suffix for all the images
  - specifying this flag when including images, you can omit the extension 
  - example: if flag has `-images-type png`, you can write `[[bulb]]` instead of `[[build.png]]`

### Changed
- updated the README markdown file
- improve argument parsing
  - when invoked without args, an attempt is made to read from standard input

### Fixed
- 🐛 Line with no leading stars causes nil pointer dereference [#2](/issues/#2) 

## [0.2.0] - 2020-09-09
### Added
- 📝 more test cases
- 🎉 new flag `images-path` to specify the base images folder
  - now when including images, you can specify just the filename
  - example: if flag has `-images-path '/Icons/AwesomFonts'`, you can write `[[bulb.png]]` instead of `[[/Icons/AwesomFonts/build.png]]`

- 🎉 new flag `lim` to specify after how many characters to wrap the text

### Changed
- updated the README markdown file

### Fixed
- 🐛 removed `go.sum` from the .gitignore file 

## [0.1.0] - 2020-09-08
- 🎉 First release!
