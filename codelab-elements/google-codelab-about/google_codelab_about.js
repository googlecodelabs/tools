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

goog.module('googlecodelabs.CodelabAbout');

const DateTimeFormat = goog.require('goog.i18n.DateTimeFormat');
const Templates = goog.require('googlecodelabs.CodelabAbout.Templates');
const soy = goog.require('goog.soy');

/** @const {string} */
const LAST_UPDATED_ATTR = 'last-updated';

/** @const {string} */
const AUTHORS_ATTR = 'authors';

/** @const {string} */
const CODELAB_TITLE_ATTR = 'codelab-title';

/**
 * @extends {HTMLElement}
 * @suppress {reportUnknownTypes}
 */
class CodelabAbout extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-about'; }

  constructor() {
    super();

    /** @private {string} */
    this.authors_ = '';

    /** @private {string} */
    this.codelabTitle_ = '';

    /** @private {boolean} */
    this.hasSetup_ = false;

    /** @private {string} */
    this.lastUpdated_ = '';
  }

  /**
   * @export
   * @override
   */
  connectedCallback() {
    if (!this.hasSetup_) {
      this.setupDom_();
    }
  }

  /**
   * @export
   * @override
   */
  disconnectedCallback() {}

  /**
   * @return {!Array<string>}
   * @export
   */
  static get observedAttributes() {
    return [AUTHORS_ATTR, LAST_UPDATED_ATTR, CODELAB_TITLE_ATTR];
  }

  /**
   * @param {string} attr
   * @param {?string} oldValue
   * @param {?string} newValue
   * @param {?string} namespace
   * @export
   * @override
   */
  attributeChangedCallback(attr, oldValue, newValue, namespace) {
    switch(attr) {
      case LAST_UPDATED_ATTR:
        if (this.hasAttribute(LAST_UPDATED_ATTR)) {
          this.lastUpdated_ = this.getAttribute(LAST_UPDATED_ATTR);
        }
        break;
      case AUTHORS_ATTR:
        if (this.hasAttribute(AUTHORS_ATTR)) {
          this.authors_ = this.getAttribute(AUTHORS_ATTR);
        }
        break;
      case CODELAB_TITLE_ATTR:
        if (this.hasAttribute(CODELAB_TITLE_ATTR)) {
          this.codelabTitle_ = this.getAttribute(CODELAB_TITLE_ATTR);
        }
        break;
    }

    this.setupDom_();
  }

  /**
   * @private
   * @param {?string} dateString
   * @return {?string}
   */
  static formatDate_(dateString) {
    if (!dateString) {
      return null;
    }
    // Formatting the Last updated date.
    const lastUpdatedDate = new Date(dateString);
    const dateFormat = new DateTimeFormat('MMM d, yyyy');
    return dateFormat.format(lastUpdatedDate);
  }

  /**
   * @private
   */
  setupDom_() {
    // Generate the content using a soy template.
    soy.renderElement(this, Templates.about, {
      lastUpdated: CodelabAbout.formatDate_(this.lastUpdated_),
      authors: this.authors_,
      codelabTitle: this.codelabTitle_.split(':').join(':||').split('||'),
    });

    this.hasSetup_ = true;
  }
}

exports = CodelabAbout;
