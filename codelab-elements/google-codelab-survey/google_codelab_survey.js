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

goog.module('googlecodelabs.CodelabSurvey');

const EventHandler = goog.require('goog.events.EventHandler');
const HTML5LocalStorage =
    goog.require('goog.storage.mechanism.HTML5LocalStorage');
const Templates = goog.require('googlecodelabs.CodelabSurvey.Templates');
const dom = goog.require('goog.dom');
const events = goog.require('goog.events');
const soy = goog.require('goog.soy');


/**
 * The prefix for all survey keys in local storage.
 * @const {string}
 */
const STORAGE_KEY_PREFIX = 'codelab-survey-';

/**
 * The id for the current survey.
 * @const {string}
 */
const SURVEY_ID_ATTR = 'survey-id';


/**
 * The upgraded id (to prevent FUOC).
 * @const {string}
 */
const SURVEY_UPGRADED_ATTR = 'upgraded';


/** @const {string} */
const DEFAULT_SURVEY_NAME = 'default-codelabs-survey';

/** @const {string} */
const OPTION_WRAPPER_CLASS = 'survey-option-wrapper';

/** @const {string} */
const RADIO_TEXT_CLASS = 'option-text';

/**
 * @extends {HTMLElement}
 */
class CodelabSurvey extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-survey'; }

  constructor() {
    super();

    /**
     * The name of the survey
     * @private {string}
     * @const
     */
    this.surveyName_ = this.getAttribute(SURVEY_ID_ATTR) || DEFAULT_SURVEY_NAME;

    /**
     * @private {!HTML5LocalStorage}
     * @const
     */
    this.storage_ = new HTML5LocalStorage();

    /**
     * @private {string}
     * @const
     */
    this.storageKey_ = STORAGE_KEY_PREFIX + this.surveyName_;

    /**
     * @private {!Object<string, !Object>}
     * @const
     */
    this.storedData_ = {};

    /**
     * @private {!EventHandler}
     * @const
     */
    this.eventHandler_ = new EventHandler();
  }

  /**
   * @export
   * @override
   */
  connectedCallback() {
    this.checkStoredData_();
    this.updateDom_();
    this.bindEvents_();
  }

  /** @private */
  bindEvents_() {
    this.eventHandler_.listen(this, events.EventType.CHANGE,
      (event) => this.handleOptionSelected_(event));
  }

  /**
   * @param {!Event} event
   * @private
   */
  handleOptionSelected_(event) {
    if (!(event.target instanceof HTMLInputElement)) {
      return;
    }
    const inputElement = event.target;
    const optionWrapperElement =
        dom.getAncestorByClass(inputElement, OPTION_WRAPPER_CLASS);
    if (!(optionWrapperElement instanceof Element)) {
      return;
    }
    const optionTextElement =
        optionWrapperElement.querySelector(`.${RADIO_TEXT_CLASS}`);
    let answer = '';
    if (optionTextElement) {
      answer = optionTextElement.textContent;
    }

    const question = inputElement.name;
    this.storedData_[this.surveyName_][question] = answer;
    this.storage_.set(
      this.storageKey_, JSON.stringify(this.storedData_[this.surveyName_]));
    const codelabEvent = new CustomEvent('google-codelab-action', {
      detail: {
        'category': 'survey',
        'action': question.substring(0, 500),
        'label': answer.substring(0, 500)
      }
    });
    document.body.dispatchEvent(codelabEvent);
  }

  /** @private */
  checkStoredData_() {
    const storedData = this.storage_.get(this.storageKey_);
    if (storedData) {
      this.storedData_[this.surveyName_] = /** @type {!Object} */ (
        JSON.parse(storedData));
    } else {
      this.storedData_[this.surveyName_] = {};
    }
  }

  /** @private */
  updateDom_() {
    const radioGroupEls = this.querySelectorAll('paper-radio-group');
    const questionEls = this.querySelectorAll('h4');
    const surveyQuestions = [];
    if (radioGroupEls.length && (questionEls.length == radioGroupEls.length)) {
      radioGroupEls.forEach((radioGroupEl, index) => {
        const surveyOptions = [];
        const polymerRadioEls = radioGroupEl.querySelectorAll(
          'paper-radio-button');
        dom.removeNode(radioGroupEl);
        polymerRadioEls.forEach(radioEl => {
          const title = radioEl.textContent;
          surveyOptions.push({
            radioId: this.normalizeIdAttr_(questionEls[index].textContent, title),
            radioTitle: title
          });
        });
        surveyQuestions.push({
          question: questionEls[index].textContent,
          options: surveyOptions
        });
        dom.removeNode(questionEls[index]);
      });
      const updatedDom = soy.renderAsElement(Templates.survey, {
        surveyName: this.surveyName_,
        surveyQuestions: surveyQuestions
      });
      this.appendChild(updatedDom);
    }
    this.setAnsweredQuestions_();
    this.setAttribute(SURVEY_UPGRADED_ATTR, '');
  }

  /** @private */
  setAnsweredQuestions_() {
    const surveyData = this.storedData_[this.surveyName_];
    if (surveyData) {
      Object.keys(surveyData).forEach(key => {
        const id = this.normalizeIdAttr_(key, surveyData[key]);
        /** @type {?HTMLInputElement} */
        const inp = /** @type {?HTMLInputElement} */ (
            this.querySelector(`#${id}`));
        if (inp) {
          inp.checked = true;
        }
      });
    }
  }

  /**
   * @param {string} question
   * @param {string} answer
   * @return {string}
   * @private
   */
  normalizeIdAttr_(question, answer) {
    return `${question}--${answer}`.replace(/\s+/g, '-')
        .replace(/[^a-zA-Z0-9 \-]/g, '').toLowerCase();
  }

  /**
   * @export
   * @override
   */
  disconnectedCallback() {
    this.eventHandler_.removeAll();
  }
}

exports = CodelabSurvey;
