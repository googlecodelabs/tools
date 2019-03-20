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

goog.module('googlecodelabs.CodelabIndex.Cards');

const HTML5LocalStorage = goog.require('goog.storage.mechanism.HTML5LocalStorage');
const Templates = goog.require('googlecodelabs.CodelabIndex.Templates');
const dom = goog.require('goog.dom');
const soy = goog.require('goog.soy');

/** @const {string} */
const AUTHORS_ATTR = 'authors';

/** @const {string} */
const CATEGORY_ATTR = 'category';

/** @const {string} */
const CATEGORY_PARAM = 'cat';

/** @const {string} */
const TITLE_ATTR = 'title';

/** @const {string} */
const DURATION_ATTR = 'duration';

/** @const {string} */
const UPDATED_ATTR = 'updated';

/** @const {string} */
const TAGS_ATTR = 'tags';

/** @const {string} */
const SORT_ATTR = 'sort';

/** @const {string} */
const FILTER_ATTR = 'filter';

/** @const {string} */
const SORT_ALPHA = 'alpha';

/** @const {string} */
const SORT_RECENT = 'recent';

/** @const {string} */
const SORT_DURATION = 'duration';

/** @const {string} */
const HIDDEN_ATTR = 'hidden';

/** @const {string} */
const PROGRESS_ATTR = 'progress';

/** @const {string} */
const STEPS_ATTR = 'steps';

/** @const {string} */
const NUM_ATTR = 'num';

/**
 * @extends {HTMLElement}
 */
class Cards extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-index-cards'; }

  constructor() {
    super();

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
    if (!this.hasAttribute(SORT_ATTR)) {
      this.setAttribute(SORT_ATTR, SORT_ALPHA);
    } else {
      this.sort_();
    }

    if (this.hasAttribute(FILTER_ATTR) ||
        this.hasAttribute(CATEGORY_ATTR) ||
        this.hasAttribute(TAGS_ATTR)) {
      this.filter_();
    }
  }

  /**
   * @return {!Array<string>}
   * @export
   */
  static get observedAttributes() {
    return [SORT_ATTR, FILTER_ATTR, CATEGORY_ATTR, TAGS_ATTR];
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
      case SORT_ATTR:
        this.sort_();
        break;
      case FILTER_ATTR:
      case CATEGORY_ATTR:
      case TAGS_ATTR:
        this.filter_();
        break;
    }
  }


  /**
   * @private
   */
  sort_() {
    let sort = this.getAttribute(SORT_ATTR) || SORT_ALPHA;
    const cards = [...this.querySelectorAll('.card')];
    if (cards.length < 2) {
      // No point sorting 0 or 1 items.
      return;
    }

    switch (sort) {
      case SORT_DURATION:
        cards.sort(this.sortDuration_.bind(this));
        break;
      case SORT_RECENT:
        cards.sort(this.sortRecent_.bind(this));
        break;
      case SORT_ALPHA:
      default:
        sort = SORT_ALPHA;
        cards.sort(this.sortAlpha_.bind(this));
        break;
    }

    cards.forEach((card) => this.appendChild(card));

    const url = new URL(document.location.toString());
    if (!sort || sort === SORT_ALPHA) {
      url.searchParams.delete(SORT_ATTR);
    } else {
      url.searchParams.set(SORT_ATTR, sort);
    }

    this.setNum_();

    const path = `${url.pathname}${url.search}`;
    window.history.replaceState({path}, document.title, path);
  }

  /**
   * @private
   */
  setNum_() {
    const num = this.querySelectorAll(`.card:not([${HIDDEN_ATTR}])`).length;
    this.setAttribute(NUM_ATTR, num);
  }

  /**
   * @param {!Element} a
   * @param {!Element} b
   * @return {number}
   * @private
   */
  sortDuration_(a, b) {
    if (!a || !b) {
      return 0;
    }
    const aDuration = parseFloat(a.getAttribute(DURATION_ATTR)) || 0;
    const bDuration = parseFloat(b.getAttribute(DURATION_ATTR)) || 0;
    const diff = aDuration - bDuration;
    if (diff === 0) {
      return this.sortRecent_(a, b);
    } else {
      return diff;
    }
  }

  /**
   * @param {!Element} a
   * @param {!Element} b
   * @return {number}
   * @private
   */
  sortRecent_(a, b) {
    if (!a || !b) {
      return 0;
    }
    const aUpdated = new Date(a.getAttribute(UPDATED_ATTR) || 0);
    const bUpdated = new Date(b.getAttribute(UPDATED_ATTR) || 0);
    const diff = bUpdated.getTime() - aUpdated.getTime();
    if (diff === 0) {
      return this.sortAlpha_(a, b);
    } else {
      return diff;
    }
  }

  /**
   * @param {!Element} a
   * @param {!Element} b
   * @return {number}
   * @private
   */
  sortAlpha_(a, b) {
    if (!a || !b) {
      return 0;
    }
    const aTitle = a.getAttribute(TITLE_ATTR);
    const bTitle = b.getAttribute(TITLE_ATTR);
    if (aTitle < bTitle) {
      return -1;
    } else if (aTitle > bTitle) {
      return 1;
    } else {
      return 0;
    }
  }

  /**
   * @private
   */
  filter_() {
    const filter = this.normalizeValue_(this.getAttribute(FILTER_ATTR));
    const tags = this.cleanStrings_(
        (this.getAttribute(TAGS_ATTR) || '').split(','));
    const categories = this.cleanStrings_(
        (this.getAttribute(CATEGORY_ATTR) || '').split(','));

    const cards = [...this.querySelectorAll('.card')];
    cards.forEach((card) => {
      const title = this.normalizeValue_(card.getAttribute(TITLE_ATTR));
      const cardCategories = this.cleanStrings_(
          (card.getAttribute(CATEGORY_ATTR) || '').split(','));
      const cardTags = this.cleanStrings_(
          (card.getAttribute(TAGS_ATTR) || '').split(','));

      let matchesFilter = true;
      let matchesTags = true;
      let matchesCategory = true;

      if (filter) {
        matchesFilter = title.indexOf(filter) !== -1;
      }

      if (tags.length) {
        matchesTags = this.arrayContains_(cardTags, tags);
      }

      if (categories.length) {
        matchesCategory = this.arrayContains_(cardCategories, categories);
      }

      if (matchesFilter && matchesTags && matchesCategory) {
        card.removeAttribute(HIDDEN_ATTR);
      } else {
        card.setAttribute(HIDDEN_ATTR, '');
      }
    });

    const url = new URL(document.location.toString());
    if (tags.length) {
      url.searchParams.set(TAGS_ATTR, tags.join(','));
    } else {
      url.searchParams.delete(TAGS_ATTR);
    }

    if (categories.length) {
      url.searchParams.set(CATEGORY_PARAM, categories.join(','));
    } else {
      url.searchParams.delete(CATEGORY_PARAM);
    }

    if (filter) {
      url.searchParams.set(FILTER_ATTR, filter);
    } else {
      url.searchParams.delete(FILTER_ATTR);
    }

    this.setNum_();

    const path = `${url.pathname}${url.search}`;
    window.history.replaceState({path}, document.title, path);
  }

  /**
   * Returns true if any of the items in A are in B.
   * @param {!Array<string>} a
   * @param {!Array<string>} b
   * @return {boolean}
   * @private
   */
  arrayContains_(a, b) {
    for (let i = 0; i < a.length; i++) {
      if (b.includes(a[i])) {
        return true;
      }
    }
    return false;
  }

  /**
   * @param {string} category 
   * @return {string}
   * @private
   */
  normalizeCategory_(category) {
    return category.toLowerCase()
        .replace(/\s+/g, '-')           // Replace spaces with -
        .replace(/--+/g, '-')         // Replace multiple - with single -
        .trim().split(',').shift();
  }

  /**
   * Trims whitespace and converts to lower case.
   * @param {string|undefined} v
   * @return {string}
   * @private
   */
  normalizeValue_(v) {
    return (v || '').trim().toLowerCase()
        .replace('\n', '')
        .replace(/\s+/g, ' ');
  }

  /**
   * @param {!Array<string>} strings 
   * @return {!Array<string>}
   * @private
   */
  cleanStrings_(strings) {
    strings = strings || [];
    const a = [];
    strings.forEach((s) => {
      const v = this.normalizeValue_(s);
      if (v) {
        a.push(v);
      }
    });
    return a.sort();
  }

  /**
   * @param {!HTMLAnchorElement} link
   * @export
   */
  addCard(link) {
    const info = {
      category: this.normalizeCategory_(link.getAttribute(CATEGORY_ATTR) || ''),
      title: link.getAttribute(TITLE_ATTR) || '',
      duration: parseInt(link.getAttribute(DURATION_ATTR), 10) || 0,
      updated: this.prettyDate_(link.getAttribute(UPDATED_ATTR)) || '',
      tags: link.getAttribute(TAGS_ATTR) || '',
      authors: link.getAttribute(AUTHORS_ATTR) || ''
    };
    soy.renderElement(link, Templates.card, info);
    link.classList.add('card');
    this.addHomeLinkForCard_(link);
    this.showProgressForCard_(link);
    this.appendChild(link);
  }

  /**
   * @param {!HTMLAnchorElement} link
   * @private
   */
  addHomeLinkForCard_(link) {
    const url = new URL(link.href, document.location.origin);
    if (!url.searchParams.has('index')) {
      url.searchParams.set('index', document.location.pathname);
    }
    dom.safe.setAnchorHref(link, url.href);
  }

  /**
   * @param {!Element} link
   * @private
   */
  showProgressForCard_(link) {
    const id = link.getAttribute('id');
    if (id) {
      const progress = this.storage_.get(`progress_${id}`);
      const steps = link.getAttribute(STEPS_ATTR);
      if (progress && steps) {
        link.setAttribute(PROGRESS_ATTR,
            (parseFloat(progress) / parseFloat(steps) - 1).toFixed(2));
      }
    }
  }

  /**
   * @param {string} updated
   * @returns {string}
   * @private
   */
  prettyDate_(updated) {
    if (!updated) {
      return '';
    }
    const mNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug',
        'Sep', 'Oct', 'Nov', 'Dec'];
    const d = new Date(updated);
    return mNames[d.getMonth()] + ' ' + d.getUTCDate() + ', ' + d.getFullYear();
  };
}

exports = Cards;
