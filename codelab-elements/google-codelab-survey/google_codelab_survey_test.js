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

goog.module('googlecodelabs.CodelabSurveyTest');
goog.setTestOnly();

const CodelabSurvey = goog.require('googlecodelabs.CodelabSurvey');
window.customElements.define(CodelabSurvey.getTagName(), CodelabSurvey);
const HTML5LocalStorage =
    goog.require('goog.storage.mechanism.HTML5LocalStorage');
const testSuite = goog.require('goog.testing.testSuite');
goog.require('goog.testing.asserts');
goog.require('goog.testing.jsunit');

let div;

/** @const {!HTML5LocalStorage} */
const localStorage = new HTML5LocalStorage();

const polymerHtml = '<google-codelab-survey survey-id="test">' +
  '<h4>Question?</h4><paper-radio-group>' +
  '<paper-radio-button>Title Text</paper-radio-button>' +
  '<paper-radio-button>Second Option</paper-radio-button>' +
  '</paper-radio-group></google-codelab-survey>';

const polymerHtmlInvalid = '<google-codelab-survey survey-id="test"><paper-radio-group>' +
  '<paper-radio-button>Title Text</paper-radio-button>' +
  '</paper-radio-group></google-codelab-survey>';

testSuite({

  setUp() {
    if (localStorage.isAvailable()) {
      localStorage.clear();
    }
    div = document.createElement('div');
    div.innerHTML = polymerHtml;
  },

  tearDown() {
    if (localStorage.isAvailable()) {
      localStorage.clear();
    }
    document.body.innerHTML = '';
    div = null;
  },

  testCodelabSurveyUpgraded() {
    document.body.appendChild(div);
    const surveyCE = div.querySelector('google-codelab-survey');
    const radioInputEl = surveyCE.querySelector('input#question--title-text');
    const radioLabelEl = surveyCE.querySelector('label#question--title-text-label');
    const radioTextEl = surveyCE.querySelector('.option-text');
    const surveyWrapperEl = surveyCE.querySelector('.survey-questions');
    assertNotNull(radioInputEl);
    assertEquals('Question?', radioInputEl.name);
    assertNotNull(radioLabelEl);
    assertEquals('test', surveyWrapperEl.getAttribute('survey-name', ''));
    assertEquals('Title Text', radioTextEl.textContent);
    assertEquals('question--title-text', radioLabelEl.getAttribute('for', ''));
    assertTrue(surveyCE.hasAttribute('upgraded'));
  },

  testCodelabSurveyIncorrectFormatNotUpgraded() {
    div.innerHTML = polymerHtmlInvalid;
    document.body.appendChild(div);
    const radioInputEl = div.querySelector('input#title-text');
    const radioLabelEl = div.querySelector('label#title-text-label');
    assertNull(radioInputEl);
    assertNull(radioLabelEl);
  },

  testCodelabSurveyOptionClick() {
    document.body.appendChild(div);
    const optionEls = div.querySelectorAll('.survey-option-wrapper');
    // If nothing is in local storage no options should be set.
    assertFalse(optionEls[0].querySelector('input').checked);
    assertFalse(optionEls[1].querySelector('input').checked);

    optionEls[0].click();
    assertEquals('{"Question?":"Title Text"}', localStorage.get('codelab-survey-test'));
    optionEls[1].click();
    assertEquals('{"Question?":"Second Option"}', localStorage.get('codelab-survey-test'));
  },

  testCodelabSurveyLoadsStoredAnswers() {
    localStorage.set('codelab-survey-test', '{"Question?":"Second Option"}');
    document.body.appendChild(div);
    const optionEls = div.querySelectorAll('.survey-option-wrapper');

    // Second option should be selected (answer loaded from local storage)
    assertFalse(optionEls[0].querySelector('input').checked);
    assertTrue(optionEls[1].querySelector('input').checked);
  },
});
