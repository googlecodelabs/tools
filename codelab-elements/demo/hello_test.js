/**
 * @license
 * Copyright 2018 Google Inc.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *      http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

goog.module('googlecodelabs.hello_test');
goog.setTestOnly();

const HelloElement = goog.require('googlecodelabs.HelloElement');
const testSuite = goog.require('goog.testing.testSuite');
goog.require('goog.testing.asserts');
goog.require('goog.testing.jsunit');

testSuite({
  testHelloEquals() {
    const x = 6;
    assertEquals(6, x);
  },

  testHelloUpgraded() {
    const div = document.createElement('div');
    div.innerHTML = "<hello-element>static</hello-element>";
    document.body.appendChild(div);
    let text = div.textContent;
    assert(`"${text}" does not end with 'upgraded!'`, text.endsWith("upgraded!"));
  },
});
