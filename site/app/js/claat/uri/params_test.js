/*
 * Copyright 2016 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

goog.provide('claat.uri.paramsTest');
goog.setTestOnly('claat.uri.paramsTest');

goog.require('claat.uri.params');
goog.require('goog.testing.jsunit');

function testDecodeEmpty() {
  var obj = claat.uri.params.decode('');
  assertObjectEquals(obj, {});
}

function testDecodeEmptyValue() {
  var obj = claat.uri.params.decode('one=');
  assertObjectEquals(obj, {'one': ['']});
  obj = claat.uri.params.decode('one&two=');
  assertObjectEquals(obj, {'one': [''], 'two': ['']});
}

function testDecodeSingle() {
  var obj = claat.uri.params.decode('one=1&two=2');
  assertObjectEquals(obj, {'one': ['1'], 'two': ['2']});
}

function testDecodeMulti() {
  var obj = claat.uri.params.decode('one=1&two=2&one=11');
  assertObjectEquals(obj, {'one': ['1','11'], 'two': ['2']});
}
