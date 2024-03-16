%global commit d5a31912eb9f69ea1c8fed59811089ff7c4ccebf
%global shortcommit %(echo %{commit} | cut -c1-7)
%global commitdate 19701230

Name:           macro
Release:        1%{?dist}
Summary:        Macro test case for parsing
License:        none
Version:        1.2.3~%{commitdate}git%{shortcommit}

%description
RPM spec file with macros for testing the version parser

%files
