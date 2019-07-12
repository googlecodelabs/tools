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

goog.require('claat.ui.cards.Sorter');

window['CardSorter'] = claat.ui.cards.Sorter;

/** @export */
claat.ui.cards.Sorter.prototype.sort;
/** @export */
claat.ui.cards.Sorter.prototype.filter;
/** @export */
claat.ui.cards.Sorter.prototype.filterByCategory;
/** @export */
claat.ui.cards.Sorter.prototype.filterByText;
/** @export */
claat.ui.cards.Sorter.prototype.filterByTags;
/** @export */
claat.ui.cards.Sorter.prototype.clearFilters;

/** @type {claat.ui.cards.Filter} */
var f;
/** @export */
f.cat;
/** @export */
f.text;
/** @export */
f.tags;
/** @export */
f.kioskTags;
