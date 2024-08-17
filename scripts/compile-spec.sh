#!/bin/bash

# Currently there seems to be a bug with TypeSpec CLI where it doesn't work properly when run from a different directory.
# This script is a workaround to run the tsp compile command from the root directory of the project.
# See: https://github.com/microsoft/typespec/issues/3397.

current_dir=$(pwd)
cd typespec || exit
tsp compile . --emit @typespec/openapi3 --option "@typespec/openapi3.emitter-output-dir=${current_dir}/openapi"
cd "$current_dir" || exit
