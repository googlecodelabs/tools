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

goog.module('googlecodelabs.CodelabStepTest');
goog.setTestOnly();

const CodelabStep = goog.require('googlecodelabs.CodelabStep');
window.customElements.define(CodelabStep.getTagName(), CodelabStep);
const MockControl = goog.require('goog.testing.MockControl');
const testSuite = goog.require('goog.testing.testSuite');
goog.require('goog.testing.asserts');
goog.require('goog.testing.jsunit');

let mockControl;

/**
 * @param {string} s
 * @return {string}
 */
window['prettyPrintOne'] = (s) => {return s;};

testSuite({
  setUp() {
    mockControl = new MockControl();
  },

  tearDown() {
    mockControl.$resetAll();
    mockControl.$tearDown();
  },

  testDomIsSetUpCorrectly() {
    const codelabStep = new CodelabStep();
    codelabStep.innerHTML = '<h1>Test</h1>';

    document.body.appendChild(codelabStep);

    assertNotUndefined(codelabStep.querySelector('.instructions'));
    assertNotUndefined(codelabStep.querySelector('.inner'));
    assertNotUndefined(codelabStep.querySelector('h2.step-title'));
    assertEquals('Test', codelabStep.querySelector('h1').innerHTML);

    document.body.removeChild(codelabStep);
  },

  testCodePrettyprint() {
    const mockPrettyPrint = mockControl.createMethodMock(window, 'prettyPrintOne');
    mockPrettyPrint('Code').$returns('MockCodeTest').$once();

    mockControl.$replayAll();

    const codelabStep = new CodelabStep();
    codelabStep.innerHTML = '<h1>Testing</h1><pre><code>Code</code></pre>';
    document.body.appendChild(codelabStep);

    mockControl.$verifyAll();

    assertNotEquals(-1, codelabStep.innerHTML.indexOf('<code>MockCodeTest</code>'));

    document.body.removeChild(codelabStep);
  },

  testSnippetCopy() {
    const codelabStep = new CodelabStep();
    codelabStep.innerHTML = '<h1>Testing</h1><pre><code class="test-code">Code</code></pre>';
    document.body.appendChild(codelabStep);

    document.body.addEventListener('google-codelab-action', (e) => {
      const detail = e.detail;
      assertEquals('snippet', detail['category']);
      assertEquals('copy', detail['action']);
      assertEquals('Code', detail['label']);
    });

    const copyEvent = new ClipboardEvent('copy', {
      view: window,
      bubbles: true,
      cancelable: true
    });
    document.body.querySelector('.test-code').dispatchEvent(copyEvent);

    document.body.removeChild(codelabStep);
  },

  testUpdateTitle() {
    const codelabStep = new CodelabStep();

    document.body.appendChild(codelabStep);

    let title = codelabStep.querySelector('h2.step-title');
    assertEquals('1. ', title.textContent);

    codelabStep.setAttribute('step', '3');
    title = codelabStep.querySelector('h2.step-title');
    assertEquals('4. ', title.textContent);

    codelabStep.setAttribute('label', 'test label');
    title = codelabStep.querySelector('h2.step-title');
    assertEquals('4. test label', title.textContent);

    document.body.removeChild(codelabStep);
  }
});
