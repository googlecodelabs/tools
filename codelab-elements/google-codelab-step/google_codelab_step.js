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

goog.module('googlecodelabs.CodelabStep');

const EventHandler = goog.require('goog.events.EventHandler');
const HtmlSanitizer = goog.require('goog.html.sanitizer.HtmlSanitizer');
const Templates = goog.require('googlecodelabs.CodelabStep.Templates');
const dom = goog.require('goog.dom');
const safe = goog.require('goog.dom.safe');
const soy = goog.require('goog.soy');
const {identity} = goog.require('goog.functions');

/** @const {string} */
const LABEL_ATTR = 'label';

/** @const {string} */
const STEP_ATTR = 'step';

/**
 * The general codelab action event fired for trackable interactions.
 * @const {string}
 */
const CODELAB_ACTION_EVENT = 'google-codelab-action';

/**
 * @extends {HTMLElement}
 * @suppress {reportUnknownTypes}
 */
class CodelabStep extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-step'; }

  constructor() {
    super();

    /**
     * @private {?Element}
     */
    this.instructions_ = null;

    /**
     * @private {?Element}
     */
    this.inner_ = null;

    /** @private {boolean} */
    this.hasSetup_ = false;

    /**
     * @private {number}
     */
    this.step_ = 0;

    /**
     * @private {string}
     */
    this.label_ = '';

    /**
     * @private {?Element}
     */
    this.title_ = null;

    /**
     * @private {?Element}
     */
    this.about_ = null;

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
    this.setupDom_();
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
    return [LABEL_ATTR, STEP_ATTR];
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
    if (attr === LABEL_ATTR || attr === STEP_ATTR) {
      this.updateTitle_();
    }
  }

  /**
   * @private
   */
  updateTitle_() {
    if (this.hasAttribute(LABEL_ATTR)) {
      this.label_ = this.getAttribute(LABEL_ATTR);
    }

    if (this.hasAttribute(STEP_ATTR)) {
      this.step_ = parseInt(this.getAttribute(STEP_ATTR) || '', 10);
    }

    if (!this.title_) {
      return;
    }

    const title = soy.renderAsElement(Templates.title, {
      step: this.step_,
      label: this.label_
    });

    dom.replaceNode(title, this.title_);
    this.title_ = title;
  }

  /**
   * @private
   */
  setupDom_() {
    if (this.hasSetup_) {
      return;
    }

    this.setAttribute('tabindex', '-1');

    // If there is an google-codelab-about element we keep it aside.
    const aboutElements = this.getElementsByTagName('google-codelab-about');
    if (aboutElements.length > 0) {
      this.about_ = aboutElements[0];
      this.about_.parentNode.removeChild(this.about_);
    }

    // Encapsulate instructions inside containers.
    this.instructions_ = dom.createElement('div');
    this.instructions_.classList.add('instructions');
    this.inner_ = dom.createElement('div');
    this.inner_.classList.add('inner');
    this.inner_.innerHTML = this.innerHTML;
    dom.appendChild(this.instructions_, this.inner_);
    dom.removeChildren(this);

    // Get the rendered title.
    let title = this.inner_.querySelector('.step-title');
    if (!title) {
      // Generate the title using a soy template.
      title = soy.renderAsElement(Templates.title, {
        step: this.step_,
        label: this.label_,
      });
    }
    this.title_ = title;

    // Inject the title in the containers.
    dom.insertChildAt(this.inner_, title, 0);

    // Add prettyprint to code blocks.
    const codeElements = this.inner_.querySelectorAll('pre code');
    codeElements.forEach((el) => {
      if (window['prettyPrintOne'] instanceof Function) {
        const code = window['prettyPrintOne'](el.innerHTML);
        // Sanitizer that preserves class names for syntax highlighting.
        const sanitizer =
            new HtmlSanitizer.Builder().withCustomTokenPolicy(identity).build();
        safe.setInnerHtml(el, sanitizer.sanitize(code));
      } else {
        el.classList.add('prettyprint');
      }
      this.eventHandler_.listen(
        el, 'copy', () => this.handleSnippetCopy_(el));
    });

    // Re-insert the about element before the instructions.
    if (this.about_) {
      dom.appendChild(this, this.about_);
    }
    // Insert instructions container.
    dom.appendChild(this, this.instructions_);

    this.hasSetup_ = true;
  }

  /**
   * @param {!Element} el The element on which we added this event listener.
   *     This is not the same as the target of the event, because the event
   *     target can be a child of this element.
   * @private
   */
  handleSnippetCopy_(el) {
    const event = new CustomEvent(CODELAB_ACTION_EVENT, {
      detail: {
        'category': 'snippet',
        'action': 'copy',
        'label': el.textContent.substring(0, 500)
      }
    });
    document.body.dispatchEvent(event);
  }
}

exports = CodelabStep;
