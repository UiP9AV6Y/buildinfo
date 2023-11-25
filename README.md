# buildinfo

Golang utility to generate build information suitable for embedding in code.

Applications usually expose their version information in one way or another.
Up until this point, in the context of Golang, this information was embedded into
the compiled artifact via linker flags during building. The actual information
was calculated via a series of (shell-) commands. Since all of this requires
orchestration, this was usually done using some kind of build system (Make, Bazel, ...)

`buildinfo` aims to remove all that cruft, by collecting and exposing the build
information in a single call.

Despite the name, it has actually no relation and very little similarity to
[debug.Buildinfo](https://pkg.go.dev/runtime/debug#BuildInfo),
which is available since Golang 1.18.

## Usage

Typically `buildinfo` is used as [Golang generator](https://go.dev/blog/generate).

It can create boilerplate code for setting up the required code on its own:

  `buildinfo --generate golang-embed --filename version.go`  

This will create a file named **version.go** in the current working directory.

To prevent muddying your repository, it is recommended to add the generator
output to the VCS ignore list. To allow your project to run and compile in
a development setting, where version information is not of the utmost
importance, a minimalistic version should be kept under version control
nontheless. To be more precise, create two files in the directory next
to your code with the `go:generate` instruction:

### `.gitignore`
```
buildinfo.json
```

### `buildinfo.json`
```json
{}
```

## Building

`buildinfo` (the library) does not require any pre-processing.
`buildinfo` (the tool) lives in the `tools` subdirectory. Generated
files are shipped with the source code and can be recreated using
the respective utilities.

### Man pages

```sh
cd ./tools
go install github.com/mmarkdown/mmark/v2@latest
go generate ./man
```

### Build information

```sh
cd ./tools
go install github.com/UiP9AV6Y/buildinfo/tools/cmd/buildinfo@latest
go generate ./version
```

### `buildinfo`

```sh
cd ./tools
go build -o buildinfo ./cmd/buildinfo
```

## Example

`buildinfo` dogfeeds its own product in [./tools/version/version.go](./tools/version/version.go).

## Alternatives

* [go-pogo/buildinfo](https://github.com/go-pogo/buildinfo)
* [jfrog/build-info-go](https://github.com/jfrog/build-info-go)
* [hlandau/buildinfo](https://github.com/hlandau/buildinfo)

## License

BSD 4-Clause, see [LICENSE](LICENSE).
