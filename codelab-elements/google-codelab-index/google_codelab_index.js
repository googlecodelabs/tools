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

goog.module('googlecodelabs.CodelabIndex');

const Cards = goog.require('googlecodelabs.CodelabIndex.Cards');
const Debouncer = goog.require('goog.async.Debouncer');
const EventHandler = goog.require('goog.events.EventHandler');
const Templates = goog.require('googlecodelabs.CodelabIndex.Templates');
const dom = goog.require('goog.dom');
const events = goog.require('goog.events');
const soy = goog.require('goog.soy');

/** @const {string} */
const CATEGORY_ATTR = 'category';

/** @const {string} */
const CATEGORY_PARAM = 'cat';

/** @const {string} */
const SORT_ATTR = 'sort';

/** @const {string} */
const FILTER_ATTR = 'filter';

/** @const {string} */
const TAGS_ATTR = 'tags';

/**
 * Interval over which to debounce in ms.
 * @const {number}
 */
const SEARCH_DEBOUNCE_INTERVAL = 20;

/**
 * @extends {HTMLElement}
 */
class CodelabIndex extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-index'; }

  constructor() {
    super();

    /**
     * @private {!EventHandler}
     * @const
     */
    this.eventHandler_ = new EventHandler();

    /** @private {boolean} */
    this.hasSetup_ = false;

    /** @private {?Cards} */
    this.cards_ = null;

    /** @private {?Element} */
    this.sortBy_ = null;

    /** @private {?HTMLInputElement} */
    this.search_ = null;

    /** @private {?Element} */
    this.clearSearchBtn_ = null;

    /** @private {?HTMLSelectElement} */
    this.categoriesSelect_ = null;

    /**
     * @private {!Debouncer}
     * @const
     */
    this.searchDebouncer_ = new Debouncer(
      () => this.handleSearchDebounced_(), SEARCH_DEBOUNCE_INTERVAL);

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

    window.requestAnimationFrame(() => {
      document.body.removeAttribute('unresolved');
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
   * @private
   */
  addEvents_() {
    if (this.sortBy_) {
      const tabs = this.sortBy_.querySelector('#sort-by-tabs');
      if (tabs) {
        this.eventHandler_.listen(tabs, events.EventType.CLICK,
          (e) => {
            e.preventDefault();
            this.handleSortByClick_(e);
          });
      }
    }

    if (this.search_) {
      this.eventHandler_.listen(this.search_, events.EventType.KEYUP,
        () => this.handleSearch_());
    }

    if (this.clearSearchBtn_) {
      this.eventHandler_.listen(this.clearSearchBtn_, events.EventType.CLICK,
        () => this.clearSearch_());
    }

    if (this.categoriesSelect_) {
      this.eventHandler_.listen(this.categoriesSelect_, events.EventType.CHANGE,
        () => this.selectCategories_());
    }
  }

  /**
   * @private
   */
  selectCategories_() {
    if (this.cards_ && this.categoriesSelect_) {
      this.cards_.setAttribute(CATEGORY_ATTR, this.categoriesSelect_.value);
    }
  }

  /**
   * @private
   */
  clearSearch_() {
    if (this.search_) {
      this.search_.value = '';
    }
    this.handleSearch_();
  }

  /**
   * @private
   */
  handleSearch_() {
    window.requestAnimationFrame(() => this.searchDebouncer_.fire());
  }

  /**
   * @private
   */
  handleSearchDebounced_() {
    const search = /** @type {!HTMLInputElement} */ (this.search_);
    const val = search.value.trim();
    if (this.clearSearchBtn_) {
      if (val === '') {
        this.clearSearchBtn_.setAttribute('hide', '');
      } else {
        this.clearSearchBtn_.removeAttribute('hide');
      }
    }

    if (this.cards_) {
      this.cards_.setAttribute(FILTER_ATTR, val);
    }
  }

  /**
   * @param {!Event} e
   * @private
   */
  handleSortByClick_(e) {
    const target = /** @type {!Element} */ (e.target);
    const sort = target.getAttribute(SORT_ATTR);
    if (this.cards_) {
      this.cards_.setAttribute(SORT_ATTR, sort);
    }
    const selected = this.querySelector('[selected]');
    if (selected) {
      selected.removeAttribute('selected');
    }
    target.setAttribute('selected', '');
  }

  /**
   * @private
   */
  setupDom_() {
    const mainInner = this.querySelector('main .main-inner');
    if (!mainInner) {
      return;
    }

    this.search_ = /** @type {?HTMLInputElement} */ (
        document.querySelector('#search-field'));
    this.clearSearchBtn_ = document.querySelector('#clear-icon');

    const list = this.querySelector('main ul');
    let cards = /** @type {!Cards} */ (
        document.createElement('google-codelab-index-cards'));

    const url = new URL(document.location.toString());
    if (url.searchParams.has(TAGS_ATTR)) {
      cards.setAttribute(TAGS_ATTR,
        url.searchParams.getAll(TAGS_ATTR).join(','));
    }

    let selectedCategory = '';
    if (url.searchParams.has(CATEGORY_PARAM)) {
      const categories = url.searchParams.getAll(CATEGORY_PARAM);
      selectedCategory = categories[0].trim().toLowerCase();
      cards.setAttribute(CATEGORY_ATTR, categories.join(','));
    }

    let sort = 'alpha';
    if (url.searchParams.has(SORT_ATTR)) {
      sort = /** @type {string} */ (url.searchParams.get(SORT_ATTR));
      cards.setAttribute(SORT_ATTR, sort);
    }

    if (list) {
      [...list.querySelectorAll('a')].forEach((link) => {
        cards['addCard'](link);
      });
      dom.removeNode(list);
      dom.appendChild(mainInner, cards);
    } else {
      cards = mainInner.querySelector('google-codelab-index-cards');
    }

    if (cards) {
      const categories = new Set();
      [...cards.querySelectorAll('.card')].forEach((card) => {
        const category = card.getAttribute(CATEGORY_ATTR);
        if (category) {
          category.split(',').forEach((c) => {
            categories.add(c.trim());
          });
        }
      });

      const sortBy = soy.renderAsElement(Templates.sortby, {
        categories: Array.from(categories).sort(),
        selectedCategory: selectedCategory,
        sort: sort
      });
      sortBy.setAttribute('id', 'sort-by');
      dom.insertSiblingBefore(sortBy, cards);

      this.sortBy_ = sortBy;
      this.cards_ = /** @type {!Cards} */ (cards);
      this.categoriesSelect_ = /** @type {?HTMLSelectElement} */ (
          this.sortBy_.querySelector('#codelab-categories'));

      if (selectedCategory && this.categoriesSelect_) {
        [...this.categoriesSelect_.options].forEach((option) => {
          if (option.value.toLowerCase() === selectedCategory) {
            option.selected = true;
          }
        });
      }
    }

    if (url.searchParams.has(FILTER_ATTR)) {
      const filter = /** @type {string} */ (url.searchParams.get(FILTER_ATTR));
      if (this.search_) {
        this.search_.value = filter;
        this.handleSearch_();
      }
    }

    this.hasSetup_ = true;
  }
}

exports = CodelabIndex;
