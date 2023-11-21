%%%
 title = "buildinfo.json 5"
 area = "BuildInfo"
 workgroup = "BuildInfo"
%%%

## Name

*buildinfo.json* - metadata describing a project.

## Description

The various build information about a project are collected in a structured file
for machine processing, formatted in JavaScript Object Notation (JSON). The content
is a flat object with the following keys:

Whitespaces and newlines are optional. Since this file is intended to be embedded
into binaries, it is recommended to reduce its size as much as possible to avoid
unecessary bloat.

## Examples

For the purpose of readability, the following examples are shown with extra whitespaces
and newlines.

The smallest possible version is simply an empty object. It is recommended to place such
a rendition under version control to always have a valid build information source available.

~~~ json
{}
~~~

A definition with all information present:

~~~ json
{
  "version": "0.1.0",
  "revision": "deadbeefcafe",
  "branch": "trunk",
  "user": "root",
  "host": "localhost",
  "date": "1970-01-01T01:02:03.123456789Z"
}
~~~

## Authors

Gordon Bleux.

## Copyright

BSD 4-Clause

## See Also

buildinfo(1).

