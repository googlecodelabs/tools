#!/usr/bin/env bash
set -euo pipefail

# TODO Implement https://github.com/googlecodelabs/tools/issues/881 and remove this! ;-)

../claat/bin/claat serve &
claatServePID=$?
trap 'kill ${claatServePID}' EXIT

echo codelab.md | entr ../claat/bin/claat export codelab.md
