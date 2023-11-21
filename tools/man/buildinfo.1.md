%%%
 title = "buildinfo 1"
 area = "BuildInfo"
 workgroup = "BuildInfo"
%%%

## BuildInfo

*buildinfo* - utility to generate build information suitable for embedding in code.

## Synopsis

*buildinfo* **[OPTION]**...

## Description

Applications usually expose their version information in one way or another.
Up until this point, in the context of Golang, this information was embedded into
the compiled artifact via linker flags during building. The actual information
was calculated via a series of (shell-) commands. Since all of this requires
orchestration, this was usually done using some kind of build system (Make, Bazel, ...)

`buildinfo` aims to remove all that cruft, by collecting and exposing the build
information in a single call.

Available options:

**--generate** **FORMAT**
: output format to generate. valid values include *buildinfo*, *golang-embed*

**--filename** **FILE**
: write rendered data to **FILE**.

**--project-dir** **DIR**
: search for VCS root in **DIR**.

**--log.level** **LEVEL**
: progress verbosity. valid values include *debug*, *info*, *warn*, *error*

**--generate.namespace** **NAME**
: use **NAME** as the package namespace for the generated code.
  only applies when using one of the following code generators:
  *golang-embed*

**--version**
: show version and quit.

## Authors

Gordon Bleux.

## Copyright

BSD 4-Clause

## See Also

buildinfo.json(5).

