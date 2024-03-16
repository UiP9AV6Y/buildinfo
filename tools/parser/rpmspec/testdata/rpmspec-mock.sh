#!/bin/sh -eu

mock_broken_spec() {
  echo "error: Release field must be present in package: (main package)" >&2
  echo "error: query of specfile $1 failed, can't parse" >&2
  exit 1
}

mock_macro_spec() {
  case "$1" in
  "%{version}") echo "1.2.3~19701230gitd5a3191" ;;
  "%{release}") echo "1.rhel" ;;
  *) echo "Unsupported spec query '$1'" >&2 ; exit 1 ;;
  esac
}

mock_minimal_spec() {
  case "$1" in
  "%{version}") echo "1.0" ;;
  "%{release}") echo "1" ;;
  *) echo "Unsupported spec query '$1'" >&2 ; exit 1 ;;
  esac
}

if test $# -lt 4; then
  echo "Not enough arguments" >&2
  exit 1
fi

if ! test -s "$4"; then
  echo "No such file: $4" >&2
  exit 1
fi

SPEC_QUERY="$3"
SPEC_FILE="$4"
MOCK_STRATEGY=$(/bin/basename -s .spec "$SPEC_FILE")

case "$MOCK_STRATEGY" in
  broken) mock_broken_spec "$SPEC_FILE" ;;
  macro) mock_macro_spec "$SPEC_QUERY" ;;
  minimal) mock_minimal_spec "$SPEC_QUERY" ;;
  multiple) mock_minimal_spec "$SPEC_QUERY" ;;
  *)
    echo "Invalid mock strategy $MOCK_STRATEGY" >&2
    exit 1
    ;;
esac

:
