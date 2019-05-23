/*
 * Copyright 2016 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

goog.provide('claat.ui.cards.Filter');
goog.provide('claat.ui.cards.Order');
goog.provide('claat.ui.cards.Sorter');

/**
 * Kiosk tags have special meaning in that the filtering result
 * should always be a non-empty intersection between tags and kioskTags.
 *
 * @typedef {{
 *   cat: (string|undefined),
 *   text: (string|undefined),
 *   tags: (!Array<string>),
 *   kioskTags: (!Array<string>|undefined)
 * }}
 */
claat.ui.cards.Filter;

/**
 * @enum {string}
 */
claat.ui.cards.Order = {
  ALPHA: 'a-z',
  DURATION: 'duration',
  RECENT: 'recent',
};


/**
 * Cards sorter for index pages.
 *
 * @param {!Element} container Card elements parent.
 * @constructor
 */
claat.ui.cards.Sorter = function(container) {
  /**
   * Cards sorting order.
   * @private {!claat.ui.cards.Order}
   */
  this.order_ = claat.ui.cards.Order.ALPHA;

  /**
   * Current filters.
   * A card matches the filters only if it matches all of the filters,
   * or all filters are empty.
   *
   * Kiosk tags have special meaning in that the filtering result
   * should always be a non-empty intersection between tags and kioskTags.
   *
   * @private {claat.ui.cards.Filter}
   */
  this.filters_ = {tags: []};

  /**
   * Card elements parent node.
   * @private {Element}
   */
  this.container_ = container;

  /**
   * Card elements.
   * @private {NodeList}
   */
  this.cards_ = container.querySelectorAll('.codelab-card');
  this.processCards_();
};

/**
 * Set order mode to o and do the sorting.
 * @param {claat.ui.cards.Order} o
 */
claat.ui.cards.Sorter.prototype.sort = function(o) {
  this.order_ = o;
  this.do_();
};

/**
 * Resets filters to f.
 * @param {claat.ui.cards.Filter} f
 */
claat.ui.cards.Sorter.prototype.filter = function(f) {
  this.filters_.cat = normalizeValue(f.cat);
  this.filters_.text = normalizeValue(f.text);
  this.filters_.tags = cleanStrings(f.tags);
  this.filters_.kioskTags = cleanStrings(f.kioskTags);
  this.do_();
};

/**
 * Set category filter.
 * @param {string} cat Category, case-insensitive.
 */
claat.ui.cards.Sorter.prototype.filterByCategory = function(cat) {
  this.filters_.cat = normalizeValue(cat);
  this.do_();
};

/**
 * Set description substring filter.
 * @param {string} text A description substring, case-insensitive.
 */
claat.ui.cards.Sorter.prototype.filterByText = function(text) {
  this.filters_.text = normalizeValue(text);
  this.do_();
};

/**
 * Set tags filter.
 * @param {Array<string>} tags Case-insensitive tags.
 * @param {Array<string>=} opt_kioskTags Case-insensitive kiosk tags.
 */
claat.ui.cards.Sorter.prototype.filterByTags = function(tags, opt_kioskTags) {
  this.filters_.tags = cleanStrings(tags);
  this.filters_.kioskTags = cleanStrings(opt_kioskTags);
  this.do_();
};

/**
 * Remove all filters.
 */
claat.ui.cards.Sorter.prototype.clearFilters = function() {
  this.filter({tags: [], kioskTags: []});
};

/**
 * Pre-compute cards properties for faster matching.
 * @private
 */
claat.ui.cards.Sorter.prototype.processCards_ = function() {
  var pin = 0;
  for (var i = 0; i < this.cards_.length; i++) {
    var card = this.cards_[i];
    // filtering
    card.desc = (card.dataset['title'] || '').trim().toLowerCase();
    card.cats = cleanStrings((card.dataset['category'] || '').split(','));
    card.tags = cleanStrings((card.dataset['tags'] || '').split(','));
    // sorting
    card.updated = new Date(card.dataset['updated']);
    card.duration = parseInt(card.dataset['duration'], 10);
    if (card.dataset['pin']) {
      pin += 1;
      card.pin = pin;
    }
  }
};

/**
 * Loops through cards_, modifying their display style property accordingly
 * to the current filters.
 * @protected
 */
claat.ui.cards.Sorter.prototype.do_ = function() {
  var elems = Array.prototype.slice.call(this.cards_, 0);
  var n = elems.length;
  while(n--) {
    if (!this.match_(elems[n])) {
      elems.splice(n, 1);
    }
  }
  this.sort_(elems);
  for (var i = 0; i < this.cards_.length; i++) {
    var c = this.cards_[i];
    if (c.parentNode) {
      c.parentNode.removeChild(c);
    }
  }
  elems.forEach(this.container_.appendChild.bind(this.container_));
};

/**
 * Performs sorting of cards_.
 * @param {Array.<Element>} cards Card elements to sort.
 * @protected
 */
claat.ui.cards.Sorter.prototype.sort_ = function(cards) {
  switch (this.order_) {
    case claat.ui.cards.Order.DURATION:
      cards.sort(function(a, b) {
        var n = comparePinned(a, b);
        if (n !== null) {
          return n;
        }
        return a.duration - b.duration;
      });
      break;
    case claat.ui.cards.Order.RECENT:
      cards.sort(function(a, b) {
        var n = comparePinned(a, b);
        if (n !== null) {
          return n;
        }
        if (b.updated < a.updated) {
          return -1;
        }
        if (b.updated > a.updated) {
          return 1;
        }
        return 0;
      });
      break;
    default:
      // alphabetical sort
      cards.sort(function(a, b) {
        var n = comparePinned(a, b);
        if (n !== null) {
          return n;
        }
        if (a.dataset['title'] < b.dataset['title']) {
          return -1;
        }
        if (a.dataset['title'] > b.dataset['title']) {
          return 1;
        }
        return 0;
      });
  }
};

/**
 * Match a card element against current filters.
 * A card always matches against empty filters.
 *
 * @param {Element} card The card to match against the filters.
 * @return {boolean} True if the card matches the filters and should be visible.
 * @protected
 */
claat.ui.cards.Sorter.prototype.match_ = function(card) {
  // Special kiosk tags match goes first, if any.
  if (this.filters_.kioskTags && this.filters_.kioskTags.length > 0) {
    if (!intersect(this.filters_.kioskTags, card.tags)) {
      return false;
    }
  }

  // category contains exact match
  // If no filter is set, we let everything through. If a filter is set, we
  // assume this card is not a match until proven otherwise.
  var catMatch = !this.filters_.cat;
  if (this.filters_.cat) {
    for (var i = 0; i < card.cats.length; i++) {
      if (this.filters_.cat === card.cats[i]) {
        catMatch = true;
      }
    }
  }
  if (!catMatch) {
    return false;
  }

  // description substring
  if (this.filters_.text && card.desc.indexOf(this.filters_.text) === -1) {
    return false;
  }

  // Both filters_.tags and card.tags must be sorted.
  if (this.filters_.tags.length > 0) {
    if (!intersect(this.filters_.tags, card.tags)) {
      return false;
    }
  }

  // all non-empty filters match
  return true;
};

/**
 * Trims whitespace and converts to lower case.
 * @param {string|undefined} v
 * @return {string}
 */
function normalizeValue(v) {
  return (v || '').trim().toLowerCase();
}

/**
 * Trims whitespace and removes empty strings.
 * @param {Array<string>|undefined} strings
 * @return {!Array<string>} Cleaned strings sorted in ascending lexical order.
 */
function cleanStrings(strings) {
    strings = strings || [];
    var a = [];
    for (var i = 0; i < strings.length; i++) {
      var v = normalizeValue(strings[i]);
      if (v) {
        a.push(v);
      }
    }
    a.sort();
    return a;
}

/**
 * @param {Element} a
 * @param {Element} b
 * @return {number|null}
 */
function comparePinned(a, b) {
  if (a.pin && !b.pin) {
    return -1;
  }
  if (!a.pin && b.pin) {
    return 1;
  }
  if (a.pin && b.pin) {
    return a.pin - b.pin;
  }
  return null;
}

/**
 * Reports whether a and b have a non-empty intersection.
 * Computes in O(min(a.length, b.length)).
 *
 * @param {!Array<string>} a Sorted array
 * @param {!Array<string>} b Sorted array
 * @return {boolean}
 */
function intersect(a, b) {
  var i = 0;
  var j = 0;
  while (i < a.length && j < b.length) {
    if (a[i] < b[j]) {
      i++;
      continue;
    }
    if (a[i] > b[j]) {
      j++;
      continue;
    }
    return true;
  }
  return false;
}
