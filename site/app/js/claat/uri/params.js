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

goog.provide('claat.uri.params');


/**
 * Decodes URL search parts.
 *
 * @param {string} search URL search part, after the '?'.
 * @return {Object<string,!Array<string>>} An object.
 */
claat.uri.params.decode = function(search) {
  var obj = /** @dict */ {};
  if (!search) {
    return obj;
  }

  var parts = search.split('&');
  while (parts.length > 0) {
    var name = goog.global.decodeURIComponent(parts.shift());
    var value = '';
    var i = name.indexOf('=');
    if (i > 0) {
      value = name.substring(i+1);
      name = name.substring(0, i);
    }
    var a = obj[name];
    if (!a) {
      a = [];
      obj[name] = a;
    }
    a.push(value);
  }

  return obj;
};
