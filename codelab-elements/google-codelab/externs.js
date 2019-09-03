/**
 * @license
 * Copyright 2019 Google Inc.
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

/**
 * @fileoverview Externs to be included if you include Google Feedback API into
 * your compiled JS code.
 * @externs
 */

/**
 * namespace
 * @suppress {duplicate}
 */
var userfeedback = {};
userfeedback.api = {};


/**
 * @param {Object} configuration Product configuration.
 * @param {Object=} opt_productData Data about the product.
 * @return {boolean} true if Feedback was loaded, false otherwise.
 */
userfeedback.api.startFeedback = function(configuration, opt_productData) {};
