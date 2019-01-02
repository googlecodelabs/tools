#!/bin/bash
# Copyright 2018 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

# Make this repo's bazel workspace the current work dir.
# Relative to this script location.
cd $(dirname $0)/../..
pwd
bazel info

BAZEL_FLAGS="--color=no \
       --curses=no \
       --verbose_failures \
       --show_task_finish \
       --show_timestamps"

# TODO(#2): Use more sensitive build/test targets when CI is working.
bazel build -s $BAZEL_FLAGS ...
bazel test -s $BAZEL_FLAGS --test_output=all --test_arg=-debug ...
