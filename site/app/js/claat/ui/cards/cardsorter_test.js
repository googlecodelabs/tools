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

goog.provide('claat.ui.cards.SorterTest');
goog.setTestOnly('claat.ui.cards.SorterTest');

goog.require('claat.ui.cards.Filter');
goog.require('claat.ui.cards.Order');
goog.require('claat.ui.cards.Sorter');

goog.require('goog.dom');
goog.require('goog.testing.jsunit');


var sandbox;
var container;
var sorter;

function childrenId() {
  var elems = container.querySelectorAll('.codelab-card');
  return Array.from(elems).map(function(el) { return el.id });
}

function setUp() {
  sandbox = goog.dom.getElement('sandbox');

  var cards = goog.dom.getElement('original');
  container = cards.cloneNode(true);
  container.id = "cards";

  sandbox.appendChild(container);
  sorter = new claat.ui.cards.Sorter(container);
}

function tearDown() {
  goog.dom.removeChildren(sandbox);
}

function testConstructorNoFilterOrOrder() {
  assertArrayEquals(['one', 'two', 'three', 'four'], childrenId());
}

function testFilterByWebTag() {
  sorter.filterByTags([' weB ']);
  assertArrayEquals(['three'], childrenId());
}

function testFilterByCommonTag() {
  sorter.filterByTags(['common']);
  var expected = [
    'two', // pinned
    'one',
  ];
  assertArrayEquals(expected, childrenId());
}

function testFilterByCategory() {
  sorter.filterByCategory(' clouD ');
  assertArrayEquals(['one'], childrenId());
}

function testFilterByCategoryMultiple() {
  sorter.filterByCategory(' beacons ');
  var expected = [
    'three',
    'four',
  ]
  assertArrayEquals(expected, childrenId());
}

function testFilterByCategoryMultiple2() {
  sorter.filterByCategory(' android ');
  var expected = [
    'two',
    'four',
  ]
  assertArrayEquals(expected, childrenId());
}

function testFilterByText() {
  sorter.filterByText('some');
  assertArrayEquals(['three'], childrenId());
}

function testFilterReset() {
  sorter.filterByTags(['nonexistent']);
  sorter.clearFilters();
  var expected = [
    'two', // pinned
    'one',
    'three',
    'four',
  ];
  assertArrayEquals(expected, childrenId());
}

function testOrderByAlpha() {
  var expected = [
    'two',   // Zzz, pinned
    'one',   // Abc
    'three', // Bcd
    'four',
  ];
  sorter.sort(claat.ui.cards.Order.ALPHA);
  assertArrayEquals(expected, childrenId());
}

function testOrderByDuration() {
  var expected = [
    'two',   // 2, pinned
    'one',   // 1
    'three', // 3
    'four', // 4
  ];
  sorter.sort(claat.ui.cards.Order.DURATION);
  assertArrayEquals(expected, childrenId());
}

function testOrderByDate() {
  var expected = [
    'two',   // 2016-06-21, pinned
    'four', //2016-06-23
    'three', // 2016-06-22
    'one',   // 2016-06-20
  ];
  sorter.sort(claat.ui.cards.Order.RECENT);
  assertArrayEquals(expected, childrenId());
}

function testFilterAndSort() {
  var expected = [
    'two', // Zzz, pinned
    'one', // Abc
  ];
  sorter.filter({tags: ['common']});
  sorter.sort(claat.ui.cards.Order.ALPHA);
  assertArrayEquals(expected, childrenId());
}

function testFilterAndSortNoPin() {
  var expected = [
    'three', // 2016-06-22
    'one',   // 2016-06-20
  ];
  sorter.filter({tags: ['web', 'cloud']});
  sorter.sort(claat.ui.cards.Order.RECENT);
  assertArrayEquals(expected, childrenId());
}
