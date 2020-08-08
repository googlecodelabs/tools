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

goog.module('googlecodelabs.Codelab');

const EventHandler = goog.require('goog.events.EventHandler');
const HTML5LocalStorage = goog.require('goog.storage.mechanism.HTML5LocalStorage');
const KeyCodes = goog.require('goog.events.KeyCodes');
const Templates = goog.require('googlecodelabs.Codelab.Templates');
const Transition = goog.require('goog.fx.css3.Transition');
const TransitionEventType = goog.require('goog.fx.Transition.EventType');
const dom = goog.require('goog.dom');
const events = goog.require('goog.events');
const soy = goog.require('goog.soy');

/**
 * Deprecated. Title causes the bowser to display a tooltip over the whole codelab.
 * Use codelab-title instead.
 * @const {string}
 */
const TITLE_ATTR = 'title';

/** @const {string} */
const CODELAB_TITLE_ATTR = 'codelab-title';

/** @const {string} */
const ENVIRONMENT_ATTR = 'environment';

/** @const {string} */
const CATEGORY_ATTR = 'category';

/** @const {string} */
const GAID_ATTR = 'codelab-gaid';

/** @const {string} */
const FEEDBACK_LINK_ATTR = 'feedback-link';

/** @const {string} */
const SELECTED_ATTR = 'selected';

/** @const {string} */
const LAST_UPDATED_ATTR = 'last-updated';

/** @const {string} */
const DURATION_ATTR = 'duration';

/** @const {string} */
const HIDDEN_ATTR = 'hidden';

/** @const {string} */
const ID_ATTR = 'id';

/** @const {string} */
const COMPLETED_ATTR = 'completed';

/** @const {string} */
const LABEL_ATTR = 'label';

/** @const {string} */
const DONT_SET_HISTORY_ATTR = 'dsh';

/** @const {string} */
const ANIMATING_ATTR = 'animating';

/** @const {string} */
const NO_TOOLBAR_ATTR = 'no-toolbar';

/** @const {string} */
const NO_ARROWS_ATTR = 'no-arrows';

/** @const {string} */
const DISAPPEAR_ATTR = 'disappear';

/** @const {number} Page transition time in seconds */
const ANIMATION_DURATION = .5;

/** @const {string} */
const DRAWER_OPEN_ATTR = 'drawer--open';

/** @const {string} */
const ANALYTICS_READY_ATTR = 'anayltics-ready';

/**
 * The general codelab action event fired for trackable interactions.
 */
const CODELAB_ACTION_EVENT = 'google-codelab-action';

/**
 * The general codelab action event fired for trackable interactions.
 */
const CODELAB_PAGEVIEW_EVENT = 'google-codelab-pageview';

/**
 * The general codelab action event fired when the codelab element is ready.
 */
const CODELAB_READY_EVENT = 'google-codelab-ready';

/**
 * @extends {HTMLElement}
 */
class Codelab extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab'; }

  constructor() {
    super();

    /** @private {?Element} */
    this.drawer_ = null;

    /** @private {?Element} */
    this.stepsContainer_ = null;

    /** @private {?NodeList} */
    this.timeContainer_ = null;

    /** @private {?Element} */
    this.titleContainer_ = null;

    /** @private {?Element} */
    this.nextStepBtn_ = null;

    /** @private {?Element} */
    this.prevStepBtn_ = null;

    /** @private {?Element} */
    this.controls_ = null;

    /** @private {?Element} */
    this.doneBtn_ = null;

    /** @private {string} */
    this.id_ = '';

    /** @private {string} */
    this.title_ = '';

    /** @private {number} */
    this.setFocusTimeoutId_ = -1;

    /** @private {!Array<!Element>} */
    this.steps_ = [];

    /** @private {number}  */
    this.currentSelectedStep_ = -1;

    /**
     * @private {!EventHandler}
     * @const
     */
    this.eventHandler_ = new EventHandler();

    /**
     * @private {!EventHandler}
     * @const
     */
    this.transitionEventHandler_ = new EventHandler();

    /** @private {boolean} */
    this.hasSetup_ = false;

    /** @private {boolean} */
    this.ready_ = false;

    /** @private {?Transition} */
    this.transitionIn_ = null;

    /** @private {?Transition} */
    this.transitionOut_ = null;

    /** @private {boolean} */
    this.resumed_ = false;

    /**
     * @private {!HTML5LocalStorage}
     * @const
     */
    this.storage_ = new HTML5LocalStorage();
  }

  /**
   * @export
   * @override
   */
  connectedCallback() {
    if (!this.hasSetup_) {
      this.setupDom_();
    }

    this.addEvents_();

    this.configureAnalytics_();
    this.showSelectedStep_();
    this.updateTitle_();
    this.toggleArrows_();
    this.toggleToolbar_();

    if (this.resumed_) {
      console.log('resumed');
      // TODO Show resume dialog
    }

    if (!this.ready_) {
      this.ready_ = true;
      this.fireEvent_(CODELAB_READY_EVENT);
      this.setAttribute(CODELAB_READY_EVENT, '');
    }
  }

  /**
   * @export
   * @override
   */
  disconnectedCallback() {
    this.eventHandler_.removeAll();
    this.transitionEventHandler_.removeAll();
  }

  /**
   * @return {!Array<string>}
   * @export
   */
  static get observedAttributes() {
    return [TITLE_ATTR, CODELAB_TITLE_ATTR, ENVIRONMENT_ATTR, CATEGORY_ATTR,
        FEEDBACK_LINK_ATTR, SELECTED_ATTR, LAST_UPDATED_ATTR, NO_TOOLBAR_ATTR,
        NO_ARROWS_ATTR, ANALYTICS_READY_ATTR];
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
    switch (attr) {
      case TITLE_ATTR:
        if (this.hasAttribute(TITLE_ATTR)) {
          this.title_ = this.getAttribute(TITLE_ATTR);
          this.removeAttribute(TITLE_ATTR);
          this.setAttribute(CODELAB_TITLE_ATTR, this.title_);
        }
        break;
      case CODELAB_TITLE_ATTR:
        this.title_ = this.getAttribute(CODELAB_TITLE_ATTR);
        this.updateTitle_();
        break;
      case SELECTED_ATTR:
        this.showSelectedStep_();
        break;
      case NO_TOOLBAR_ATTR:
        this.toggleToolbar_();
        break;
      case NO_ARROWS_ATTR:
        this.toggleArrows_();
        break;
      case ANALYTICS_READY_ATTR:
        if (this.hasAttribute(ANALYTICS_READY_ATTR)) {
          if (this.ready_) {
            this.firePageLoadEvents_();
          } else {
            this.addEventListener(CODELAB_READY_EVENT,
                () => this.firePageLoadEvents_());
          }
        }
        break;
    }
  }

  /**
   * @return {!EventHandler}
   * @export
   */
  get eventHandler() {
    return this.eventHandler_;
  }

  /**
   * @return {!Array<!Element>}
   * @export
   */
  get steps() {
    return this.steps_;
  }

  /**
   * @private
   */
  configureAnalytics_() {
    const analytics = document.querySelector('google-codelab-analytics');
    if (analytics) {
      const gaid = this.getAttribute(GAID_ATTR);
      if (gaid) {
        analytics.setAttribute(GAID_ATTR, gaid);
      }

      analytics.setAttribute(
        ENVIRONMENT_ATTR, this.getAttribute(ENVIRONMENT_ATTR));
      analytics.setAttribute(CATEGORY_ATTR, this.getAttribute(CATEGORY_ATTR));
    }
  }

  /**
   * @export
   */
  selectNext() {
    this.setAttribute(SELECTED_ATTR, this.currentSelectedStep_ + 1);
  }

  /**
   * @export
   */
  selectPrevious() {
    this.setAttribute(SELECTED_ATTR, this.currentSelectedStep_ - 1);
  }

  /**
   * @export
   * @param {number} index
   */
  select(index) {
    this.setAttribute(SELECTED_ATTR, index);
  }

  /**
   * @private
   */
  addEvents_() {
    if (this.prevStepBtn_) {
      this.eventHandler_.listen(this.prevStepBtn_, events.EventType.CLICK,
        (e) => {
          e.preventDefault();
          e.stopPropagation();
          this.selectPrevious();
        });
    }
    if (this.nextStepBtn_) {
      this.eventHandler_.listen(this.nextStepBtn_, events.EventType.CLICK,
        (e) => {
          e.preventDefault();
          e.stopPropagation();
          this.selectNext();
        });
    }

    if (this.drawer_) {
      this.eventHandler_.listen(this.drawer_, events.EventType.CLICK,
          (e) => this.handleDrawerClick_(e));

      this.eventHandler_.listen(this.drawer_, events.EventType.KEYDOWN,
            (e) => this.handleDrawerKeyDown_(e));
    }

    if (this.titleContainer_) {
      const menuBtn = this.titleContainer_.querySelector('#menu');
      if (menuBtn) {
        this.eventHandler_.listen(menuBtn, events.EventType.CLICK, (e) => {
          e.preventDefault();
          e.stopPropagation();
          if (this.hasAttribute(DRAWER_OPEN_ATTR)) {
            this.removeAttribute(DRAWER_OPEN_ATTR);
          } else {
            this.setAttribute(DRAWER_OPEN_ATTR, '');
          }
        });

        this.eventHandler_.listen(document.body, events.EventType.CLICK, (e) => {
          if (this.hasAttribute(DRAWER_OPEN_ATTR)) {
            this.removeAttribute(DRAWER_OPEN_ATTR);
          }
        });
      }
    }

    this.eventHandler_.listen(dom.getWindow(), events.EventType.POPSTATE, (e) => {
      this.handlePopStateChanged_(e);
    });

    this.eventHandler_.listen(document.body, events.EventType.KEYDOWN, (e) => {
      this.handleKeyDown_(e);
    });
  }

  /**
   * @private
   */
  toggleToolbar_() {
    if (!this.titleContainer_) {
      return;
    }

    if (this.hasAttribute(NO_TOOLBAR_ATTR)) {
      this.titleContainer_.setAttribute(HIDDEN_ATTR, '');
    } else {
      this.titleContainer_.removeAttribute(HIDDEN_ATTR);
    }
  }

  /**
   * @private
   */
  toggleArrows_() {
    if (!this.controls_) {
      return;
    }

    if (this.hasAttribute(NO_ARROWS_ATTR)) {
      this.controls_.setAttribute(HIDDEN_ATTR, '');
    } else {
      this.controls_.removeAttribute(HIDDEN_ATTR);
    }
  }

  /**
   *
   * @param {!events.BrowserEvent} e
   * @private
   */
  handleDrawerKeyDown_(e) {
    if (!this.drawer_) {
      return;
    }

    const focused = this.drawer_.querySelector(':focus');
    let li;
    if (focused) {
      li = /** @type {!Element} */ (focused.parentNode);
    } else {
      li = this.drawer_.querySelector(`[${SELECTED_ATTR}]`);
    }

    if (!li) {
      return;
    }

    let next;
    if (e.keyCode == KeyCodes.UP) {
      next = dom.getPreviousElementSibling(li);
    } else if (e.keyCode == KeyCodes.DOWN) {
      next = dom.getNextElementSibling(li);
    }

    if (next) {
      const a = next.querySelector('a');
      if (a) {
        a.focus();
      }
    }
  }

  /**
   * @param {!events.BrowserEvent} e
   * @private
   */
  handleKeyDown_(e) {
    if (e.keyCode == KeyCodes.LEFT) {
      if (document.activeElement) {
        document.activeElement.blur();
      }
      this.selectPrevious();
    } else if (e.keyCode == KeyCodes.RIGHT) {
      if (document.activeElement) {
        document.activeElement.blur();
      }
      this.selectNext();
    }
  }

  /**
   * History popState callback
   * @param {!Event} e
   * @private
   */
  handlePopStateChanged_(e) {
    if (document.location.hash) {
      this.setAttribute(DONT_SET_HISTORY_ATTR, '');
      this.setAttribute(SELECTED_ATTR, document.location.hash.substring(1));
      this.removeAttribute(DONT_SET_HISTORY_ATTR);
    }
  }

   /**
   * Updates the browser history state
   * @param {string} path The new browser state
   * @param {boolean=} replaceState optionally replace state instead of pushing
   * @export
   */
  updateHistoryState(path, replaceState=false) {
    if (replaceState) {
      window.history.replaceState({path}, document.title, path);
    } else {
      window.history.pushState({path}, document.title, path);
    }
  }

  /**
   * @param {!Event} e
   * @private
   */
  handleDrawerClick_(e) {
    let target = /** @type {!Element} */ (e.target);

    while (target !== this.drawer_) {
      if (target.tagName.toUpperCase() === 'A') {
        break;
      }
      e.preventDefault();
      e.stopPropagation();
      target = /** @type {!Element} */ (target.parentNode);
    }

    if (target === this.drawer_) {
      return;
    }

    const selected = new URL(target.getAttribute('href'), document.location.origin)
        .hash.substring(1);

    this.setAttribute(SELECTED_ATTR, selected);
  }

  /**
   * @private
   */
  updateTitle_() {
    if (!this.title_ || !this.titleContainer_) {
      return;
    }
    const newTitleEl =
        soy.renderAsElement(Templates.title, {title: this.title_});
    document.title = this.title_;
    const oldTitleEl = this.titleContainer_.querySelector('h1');
    const buttons = this.titleContainer_.querySelector('#codelab-nav-buttons');
    if (oldTitleEl) {
      dom.replaceNode(newTitleEl, oldTitleEl);
    } else {
      dom.insertSiblingAfter(newTitleEl, buttons);
    }
  }

  /**
   * @private
   */
  updateTimeRemaining_() {
    if (!this.timeContainer_ || !this.timeContainer_.length) {
      return;
    }


    let time = 0;
    for (let i = this.currentSelectedStep_; i < this.steps_.length; i++) {
      const step = /** @type {!Element} */ (this.steps_[i]);
      let n = parseInt(step.getAttribute(DURATION_ATTR), 10);
      if (n) {
        time += n;
      }
    }

    Array.prototype.forEach.call(this.timeContainer_, (timeContainer) => {
      // Hide the time container if there was no time indication.
      if (!time) {
        timeContainer.style.display = 'none';
        return;
      }

      // Update the time container with remaining time.
      const newTimeEl =  soy.renderAsElement(Templates.timeRemaining, {time});
      const oldTimeEl = timeContainer.querySelector('.time-remaining');
      if (oldTimeEl) {
        dom.replaceNode(newTimeEl, oldTimeEl);
      } else {
        dom.appendChild(timeContainer, newTimeEl);
      }
    });
  }

  /**
   * @private
   */
  setupSteps_() {
    this.steps_.forEach((step, index) => {
      step = /** @type {!Element} */ (step);
      step.setAttribute('step', index+1);
    });
  }

  /**
   * @private
   */
  showSelectedStep_() {
    // Close drawer if any.
    this.removeAttribute(DRAWER_OPEN_ATTR);

    let selected = 0;
    if (this.hasAttribute(SELECTED_ATTR)) {
      selected = parseInt(this.getAttribute(SELECTED_ATTR), 0);
    } else {
      this.setAttribute(SELECTED_ATTR, selected);
      return;
    }

    selected = Math.min(Math.max(0, parseInt(selected, 10)),
                        this.steps_.length - 1);

    if (this.currentSelectedStep_ === selected || isNaN(selected)) {
      // Either the current step is already selected or an invalid option was provided
      // do nothing and return.
      return;
    }

    const stepTitleEl = this.steps_[selected].querySelector('.step-title');
    const stepTitle = stepTitleEl ? stepTitleEl.textContent : '';
    const stepTitlePrefix = (selected + 1) + '.';
    const re = new RegExp(stepTitlePrefix, 'g');
    this.fireEvent_(CODELAB_PAGEVIEW_EVENT, {
      'page': location.pathname + '#' + selected,
      'title': stepTitle.replace(re, '').trim()
    });

    const stepToSelect = this.steps_[selected];

    if (this.currentSelectedStep_ === -1) {
      // No previous selected step, so select the correct step with no animation
      stepToSelect.setAttribute(SELECTED_ATTR, '');
    } else {
      if (this.transitionIn_) {
        this.transitionIn_.stop();
      }
      if (this.transitionOut_) {
        this.transitionOut_.stop();
      }

      this.transitionEventHandler_.removeAll();

      const transitionInInitialStyle = {};
      const transitionInFinalStyle = {
        transform: 'translate3d(0, 0, 0)'
      };

      const transitionOutInitialStyle = {
        transform: 'translate3d(0, 0, 0)'
      };
      const transitionOutFinalStyle = {};

      const currentStep = this.steps_[this.currentSelectedStep_];
      stepToSelect.setAttribute(ANIMATING_ATTR, '');

      if (this.currentSelectedStep_ < selected) {
        // Move new step in from the right
        transitionInInitialStyle['transform'] = 'translate3d(110%, 0, 0)';
        transitionOutFinalStyle['transform'] = 'translate3d(-110%, 0, 0)';
      } else {
        // Move new step in from the left
        transitionInInitialStyle['transform'] = 'translate3d(-110%, 0, 0)';
        transitionOutFinalStyle['transform'] = 'translate3d(110%, 0, 0)';
      }

      const animationProperties = [{
        property: 'transform',
        duration: ANIMATION_DURATION,
        delay: 0,
        timing: 'cubic-bezier(0.4, 0, 0.2, 1)'
      }];

      this.transitionIn_ = new Transition(stepToSelect, ANIMATION_DURATION,
          transitionInInitialStyle, transitionInFinalStyle, animationProperties);
      this.transitionOut_ = new Transition(currentStep, ANIMATION_DURATION,
        transitionOutInitialStyle, transitionOutFinalStyle, animationProperties);

      this.transitionIn_.play();
      this.transitionOut_.play();

      this.transitionEventHandler_.listenOnce(this.transitionIn_,
            [TransitionEventType.FINISH, TransitionEventType.STOP], () => {
        stepToSelect.setAttribute(SELECTED_ATTR, '');
        stepToSelect.removeAttribute(ANIMATING_ATTR);
      });

      this.transitionEventHandler_.listenOnce(this.transitionOut_,
            [TransitionEventType.FINISH, TransitionEventType.STOP], () => {
        currentStep.removeAttribute(SELECTED_ATTR);
      });
    }

    this.currentSelectedStep_ = selected;

    // Set the focus on the new step after the animation is finished becasue it
    // messes up the animation.
    clearTimeout(this.setFocusTimeoutId_);
    this.setFocusTimeoutId_ = setTimeout(() => {
      stepToSelect.focus();
    }, ANIMATION_DURATION * 1000);

    if (this.nextStepBtn_ && this.prevStepBtn_ && this.doneBtn_) {
      if (selected === 0) {
        this.prevStepBtn_.setAttribute(DISAPPEAR_ATTR, '');
      } else {
        this.prevStepBtn_.removeAttribute(DISAPPEAR_ATTR);
      }
      if (selected === this.steps_.length - 1) {
        this.nextStepBtn_.setAttribute(HIDDEN_ATTR, '');
        this.doneBtn_.removeAttribute(HIDDEN_ATTR);
        this.fireEvent_(CODELAB_ACTION_EVENT, {
          'category': 'codelab',
          'action': 'complete',
          'label': this.title_
        });
      } else {
        this.nextStepBtn_.removeAttribute(HIDDEN_ATTR);
        this.doneBtn_.setAttribute(HIDDEN_ATTR, '');
      }
    }

    if (this.drawer_) {
      const steps = this.drawer_.querySelectorAll('li');
      steps.forEach((step, i) => {
        if (i <= selected) {
          step.setAttribute(COMPLETED_ATTR, '');
        } else {
          step.removeAttribute(COMPLETED_ATTR);
        }
        if (i === selected) {
          step.setAttribute(SELECTED_ATTR, '');
        } else {
          step.removeAttribute(SELECTED_ATTR);
        }
      });
    }

    this.updateTimeRemaining_();
    if (!this.hasAttribute(DONT_SET_HISTORY_ATTR)) {
      this.updateHistoryState(`#${selected}`, true);
    }

    if (this.id_) {
      this.storage_.set(`progress_${this.id_}`,
                        String(this.currentSelectedStep_));
    }
  }

  /**
   * @private
   */
  renderDrawer_() {
    const feedback = this.getAttribute(FEEDBACK_LINK_ATTR);
    const steps = this.steps_.map((step) => step.getAttribute(LABEL_ATTR));
    soy.renderElement(this.drawer_, Templates.drawer, {steps, feedback});
    // Start Google Feedback when the feedback link is clicked, if it exists.
    const feedbackLink = this.drawer_.querySelector('#codelab-feedback');
    if (feedbackLink) {
      this.eventHandler_.listen(feedbackLink, events.EventType.CLICK,
        (e) => {
          if ('userfeedback' in window) {
            window['userfeedback']['api']['startFeedback']
                ({'productId': '5143948'});
            e.preventDefault();
          }
        });
    }
  }

  /**
   * @private
   * @return {string}
   */
  getHomeUrl_() {
    const url = new URL(document.location.toString());
    let index = url.searchParams.get('index');
    if (!index) {
      return '/';
    }

    index = index.replace(/[^a-z0-9\-]+/ig, '');
    if (!index || index.trim() === '') {
      return '/';
    }

    if (index === 'index') {
      index = '';
    }
    const u = new URL(index, document.location.origin);
    return u.pathname;
  }

  /**
   * @param {string} eventName
   * @param {!Object=} detail
   * @protected
   */
  fireEvent_(eventName, detail={}) {
    const event = new CustomEvent(eventName, {
      detail: detail,
      bubbles: true,
    });
    this.dispatchEvent(event);
  }

  /**
   * Fires events for initial page load.
   * @private
   */
  firePageLoadEvents_() {
    this.fireEvent_(CODELAB_PAGEVIEW_EVENT, {
      'page': location.pathname + '#' + this.currentSelectedStep_,
      'title': this.steps_[this.currentSelectedStep_].getAttribute(LABEL_ATTR)
    });

    window.requestAnimationFrame(() => {
      document.body.removeAttribute('unresolved');
      this.fireEvent_(CODELAB_ACTION_EVENT, {
        'category': 'codelab',
        'action': 'ready'
      });
    });
  }

  /**
   * @private
   */
  setupDom_() {
    this.steps_ = Array.from(this.querySelectorAll('google-codelab-step'));

    soy.renderElement(this, Templates.structure, {
      homeUrl: this.getHomeUrl_()
    });

    this.drawer_ = this.querySelector('#drawer');
    this.titleContainer_ = this.querySelector('#codelab-title');
    this.stepsContainer_ = this.querySelector('#steps');
    this.controls_ = this.querySelector('#controls');
    this.prevStepBtn_ = this.querySelector('#controls #previous-step');
    this.nextStepBtn_ = this.querySelector('#controls #next-step');
    this.doneBtn_ = this.querySelector('#controls #done');

    this.steps_.forEach((step) => dom.appendChild(this.stepsContainer_, step));
    this.setupSteps_();
    this.renderDrawer_();
    this.timeContainer_ = this.querySelectorAll('.codelab-time-container');

    if (document.location.hash) {
      const h = parseInt(document.location.hash.substring(1), 10);
      if (!isNaN(h) && h) {
        this.setAttribute(SELECTED_ATTR, document.location.hash.substring(1));
      }
    }

    this.id_ = this.getAttribute(ID_ATTR);
    const progress = this.storage_.get(`progress_${this.id_}`);
    if (progress && progress !== '0') {
      this.resumed_ = true;
      this.setAttribute(SELECTED_ATTR, progress);
    }

    this.hasSetup_ = true;
  }
}

exports = Codelab;
