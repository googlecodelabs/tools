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

goog.module('googlecodelabs.CodelabAnalytics');

const EventHandler = goog.require('goog.events.EventHandler');

/**
 * The general codelab action event fired for trackable interactions.
 * @const {string}
 */
const ACTION_EVENT = 'google-codelab-action';

/**
 * The general codelab pageview event fired for trackable pageviews.
 * @const {string}
 */
const PAGEVIEW_EVENT = 'google-codelab-pageview';

/**
 * The Google Analytics ID. Analytics CE will not complete initialization
 * without a valid Analytics ID value set for this.
 * @const {string}
 */
const GAID_ATTR = 'gaid';

/** @const {string} */
const CODELAB_ID_ATTR = 'codelab-id';

/**
 * The GAID defined by the current codelab.
 * @const {string}
 */
const CODELAB_GAID_ATTR = 'codelab-gaid';

/** @const {string} */
const CODELAB_ENV_ATTR = 'environment';

/** @const {string} */
const CODELAB_CATEGORY_ATTR = 'category';

/** @const {string} */
const ANALYTICS_READY_ATTR = 'anayltics-ready';

/**
 * A list of selectors whose elements are waiting for this to be set up.
 * @const {!Array<string>}
 */
const DEPENDENT_SELECTORS = ['google-codelab'];


/**
 * Event detail passed when firing ACTION_EVENT.
 *
 * @typedef {{
 *  category: string,
 *  action: string,
 *  label: (?string|undefined),
 *  value: (?number|undefined)
 * }}
 */
let AnalyticsTrackingEvent;

/**
 * Event detail passed when firing ACTION_EVENT.
 *
 * @typedef {{
 *  page: string,
 *  title: string,
 * }}
 */
let AnalyticsPageview;


/**
 * @extends {HTMLElement}
 * @suppress {reportUnknownTypes}
 */
class CodelabAnalytics extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-analytics'; }

  constructor() {
    super();

    /** @private {boolean} */
    this.hasSetup_ = false;

    /** @private {?string} */
    this.gaid_;

    /** @private {?string} */
    this.codelabId_;

    /**
     * @private {!EventHandler}
     * @const
     */
    this.eventHandler_ = new EventHandler();

    /** @private {?string} */
    this.codelabCategory_ = this.getAttribute(CODELAB_CATEGORY_ATTR) || '';

    /** @private {?string} */
    this.codelabEnv_ = this.getAttribute(CODELAB_ENV_ATTR) || '';
  }

  /**
   * @export
   * @override
   */
  connectedCallback() {
    this.gaid_ = this.getAttribute(GAID_ATTR) || '';

    if (this.hasSetup_ || !this.gaid_) {
      return;
    }

    if (!('ga' in window)) {
      this.initGAScript_().then((response) => {
        if (response) {
          this.init_();
        }
      });
    } else {
      this.init_();
    }
  }

  /** @private */
  init_() {
    this.createTrackers_();
    this.addEventListeners_();
    this.setAnalyticsReadyAttrs_();
    this.hasSetup_ = true;
  }

  /** @private */
  addEventListeners_() {
    this.eventHandler_.listen(document.body, ACTION_EVENT,
      (e) => {
        const detail = /** @type {AnalyticsTrackingEvent} */ (
          e.getBrowserEvent().detail);
        // Add tracking...
        this.trackEvent_(
          detail['category'], detail['action'], detail['label']);
      });

    this.eventHandler_.listen(document.body, PAGEVIEW_EVENT,
      (e) => {
        const detail = /** @type {AnalyticsPageview} */ (
          e.getBrowserEvent().detail);
        this.trackPageview_(detail['page'], detail['title']);
      });
  }

  /**
   * @return {!Array<string>}
   * @export
   */
  static get observedAttributes() {
    return [CODELAB_GAID_ATTR, CODELAB_ENV_ATTR, CODELAB_CATEGORY_ATTR,
            CODELAB_ID_ATTR];
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
      case GAID_ATTR:
        this.gaid_ = newValue;
        break;
      case CODELAB_GAID_ATTR:
        if (newValue && this.hasSetup_) {
          this.createCodelabGATracker_();
        }
        break;
      case CODELAB_ENV_ATTR:
        this.codelabEnv_ = newValue;
        break;
      case CODELAB_CATEGORY_ATTR:
        this.codelabCategory_ = newValue;
        break;
      case CODELAB_ID_ATTR:
        this.codelabId_ = newValue;
        break;
      default:
    }
  }

  /**
   * Fires an analytics tracking event to all configured trackers.
   * @param {string} category The event category.
   * @param {string=} opt_action The event action.
   * @param {?string=} opt_label The event label.
   * @private
   */
  trackEvent_(category, opt_action, opt_label) {
    const params = {
      // Always event for trackEvent_ method
      'hitType': 'event',
      'dimension1': this.codelabEnv_,
      'dimension2': this.codelabCategory_ || '',
      'dimension4': this.codelabId_ || undefined,
      'eventCategory': category,
      'eventAction': opt_action || '',
      'eventLabel': opt_label || '',
    };
    this.gaSend_(params);
  }

  /**
   * @param {?string=} opt_page The page to track.
   * @param {?string=} opt_title The codelabs title.
   * @private
   */
  trackPageview_(opt_page, opt_title) {
    const params = {
      'hitType': 'pageview',
      'dimension1': this.codelabEnv_,
      'dimension2': this.codelabCategory_,
      'dimension4': this.codelabId_ || undefined,
      'page': opt_page || '',
      'title': opt_title || ''
    };
    this.gaSend_(params);
  }

  /**
   * Sets analytics ready attributes on dependent elements.
   */
  setAnalyticsReadyAttrs_() {
    DEPENDENT_SELECTORS.forEach((selector) => {
      document.querySelectorAll(selector).forEach((element) => {
        element.setAttribute(ANALYTICS_READY_ATTR, ANALYTICS_READY_ATTR);
      });
    });
  }

  /** @private */
  gaSend_(params) {
    window['ga'](function() {
      if (window['ga'].getAll) {
        const trackers = window['ga'].getAll();
        trackers.forEach((tracker) => {
          tracker.send(params);
        });
      }
    });
  }

  /**
   * @export
   * @override
   */
  disconnectedCallback() {
    this.eventHandler_.removeAll();
  }

  /**
   * @return {string}
   * @private
   */
  getGAView_() {
    let parts = location.search.substring(1).split('&');
    for (let i = 0; i < parts.length; i++) {
      let param = parts[i].split('=');
      if (param[0] === 'viewga') {
        return param[1];
      }
    }
    return '';
  }

  /**
   * @return {!Promise}
   * @export
   */
  static injectGAScript() {
    /** @type {!HTMLScriptElement} */
    const resource = /** @type {!HTMLScriptElement} */ (
        document.createElement('script'));
    resource.src = 'https://www.google-analytics.com/analytics.js';
    resource.async = false;
    return new Promise((resolve, reject) => {
      resource.onload = () => resolve(resource);
      resource.onerror = (event) => {
        // remove on error
        if (resource.parentNode) {
          resource.parentNode.removeChild(resource);
        }
        reject();
      };
      if (document.head) {
        document.head.appendChild(resource);
      }
    });
  }

  /**
   * @return {!Promise}
   * @private
   */
  async initGAScript_() {
    // This is a pretty-printed version of the function(i,s,o,g,r,a,m) script
    // provided by Google Analytics.
    window['GoogleAnalyticsObject'] = 'ga';
    window['ga'] = window['ga'] || function() {
      (window['ga']['q'] = window['ga']['q'] || []).push(arguments);
    };
    window['ga']['l'] = (new Date()).valueOf();

    try {
      return await CodelabAnalytics.injectGAScript();
    } catch(e) {
      return;
    }
  }

  /** @private */
  createTrackers_() {
    if (window['ga']) {
      // The default tracker is given name 't0' per analytics.js dev docs.
      if (this.gaid_ && !this.isTrackerCreated_(this.gaid_)) {
        window['ga']('create', this.gaid_, 'auto');
      }

      const gaView = this.getGAView_();
      if (gaView && !this.isTrackerCreated_(gaView)) {
        window['ga']('create', gaView, 'auto', 'view');
        window['ga']('view.send', 'pageview');
      }
    }

    this.createCodelabGATracker_();
  }

  /**
   * Creates a GA tracker specific to the codelab.
   * @private
   */
  createCodelabGATracker_() {
    if (window['ga']) {
      const codelabGAId = this.getAttribute(CODELAB_GAID_ATTR);
      if (codelabGAId && !this.isTrackerCreated_(codelabGAId)) {
        window['ga']('create', codelabGAId, 'auto', 'codelabAccount');
      }
    }
  }

  /**
   * @param {string} trackerId The tracker ID to check for.
   * @return {boolean}
   * @private
   */
  isTrackerCreated_(trackerId) {
    let isCreated = false;
    if (window['ga'] && window['ga'].getAll) {
      const allTrackers = window['ga'].getAll();
      allTrackers.forEach((tracker) => {
        if (tracker.get('trackingId') == trackerId) {
          isCreated = true;
        }
      });
    }
    return isCreated;
  }
}

exports = CodelabAnalytics;
