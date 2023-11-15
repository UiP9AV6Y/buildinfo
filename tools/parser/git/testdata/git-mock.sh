#!/bin/sh -eu

mock_not_git_repo() {
  echo "fatal: not a git repository (or any parent up to mount point /mock)" >&2
  exit 128
}

mock_show_toplevel() {
  echo "/mock/src"
}

mock_parse() {
  case "$2" in
    "rev-parse --abbrev-ref HEAD")
      if test "$1" = "BRANCH_FAIL"; then
        echo "fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree." >&2
        exit 128
      else
        echo "test_mock"
      fi
      ;;
    "rev-parse HEAD")
      if test "$1" = "REV_FAIL"; then
        echo "fatal: ambiguous argument 'HEAD': unknown revision or path not in the working tree." >&2
        exit 128
      else
        echo "deadbeefcafe"
      fi
      ;;
    "describe --tags --abbrev=0")
      if test "$1" = "TAG_FAIL"; then
        echo "fatal: No names found, cannot describe anything." >&2
        exit 128
      else
        echo "v1.23.456"
      fi
      ;;
    *)
      echo "Invalid mock usage"; exit 1 ;;
  esac
}


if test $# -lt 2; then
  echo "Not enough arguments" >&2
  exit 1
fi

if test "$1" -ne "-C"; then
  echo "Missing working directory argument" >&2
  exit 1
fi

MOCK_STRATEGY="$2"
shift 2

case "$MOCK_STRATEGY" in
  /mock/NOT_GIT_REPO) mock_not_git_repo ;;
  /mock/SHOW_TOPLEVEL) mock_show_toplevel ;;
  /mock/PARSE_ALL) mock_parse "OK" "$*" ;;
  /mock/PARSE_TAG_FAIL) mock_parse "TAG_FAIL" "$*" ;;
  /mock/PARSE_REV_FAIL) mock_parse "REV_FAIL" "$*" ;;
  /mock/PARSE_BRANCH_FAIL) mock_parse "BRANCH_FAIL" "$*" ;;
  *)
    echo "Invalid mock strategy $MOCK_STRATEGY" >&2
    exit 1
    ;;
esac

:
