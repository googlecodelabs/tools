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

goog.module('googlecodelabs.CodelabAnalyticsTest');
goog.setTestOnly();

const CodelabAnalytics = goog.require('googlecodelabs.CodelabAnalytics');
window.customElements.define(CodelabAnalytics.getTagName(), CodelabAnalytics);
const MockControl = goog.require('goog.testing.MockControl');
const dom = goog.require('goog.dom');
const testSuite = goog.require('goog.testing.testSuite');
goog.require('goog.testing.asserts');
goog.require('goog.testing.jsunit');

let mockControl;
/**
 * Noop the inject function, we don't need to be going to the actual network
 * in tests
 */
CodelabAnalytics.injectGAscript = () => {};

testSuite({

  setUp() {
    mockControl = new MockControl();
  },

  tearDown() {
    dom.removeNode(document.body.querySelector('google-codelab-analytics'));

    mockControl.$resetAll();
    mockControl.$tearDown();
  },

  testGAIDAttr_InitsTracker() {
    const analytics = new CodelabAnalytics();

    // Need to mock as we don't have window.ga.getAll()
    const mockGetAll = mockControl.createFunctionMock('getAll');
    const mockCreate = mockControl.createFunctionMock('create');
    mockGetAll().$returns([]).$anyTimes();
    mockCreate().$once();

    analytics.setAttribute('gaid', 'UA-123');

    window['ga'] = (...args) => {
      if (['create', 'getAll'].indexOf(args[0]) !== -1) {
        window['ga'][args[0]]();
      }
    };
    window['ga']['getAll'] = mockGetAll;
    window['ga']['create'] = mockCreate;

    mockControl.$replayAll();
    document.body.appendChild(analytics);
    mockControl.$verifyAll();
  },

  testViewParam_InitsViewTracker() {
    const analytics = new CodelabAnalytics();
    analytics.setAttribute('gaid', 'UA-123');

    const loc = window.location;
    var newurl = loc.protocol + '//' + loc.host + loc.pathname +
        '?viewga=testView&param2=hi';
    window.history.pushState({ path: newurl }, '', newurl);

    // Need to mock as we don't have window.ga.getAll()
    const mockGetAll = mockControl.createFunctionMock('getAll');
    const mockCreate = mockControl.createFunctionMock('create');
    mockGetAll().$returns([]).$anyTimes();
    // Creates 2 trackers (because of view param).
    mockCreate().$times(2);

    window['ga'] = (...args) => {
      if (['create', 'getAll'].indexOf(args[0]) !== -1) {
        window['ga'][args[0]]();
      }
    };
    window['ga']['getAll'] = mockGetAll;
    window['ga']['create'] = mockCreate;

    mockControl.$replayAll();

    document.body.appendChild(analytics);

    mockControl.$verifyAll();
    window.history.back();
  },

  testCodelabGAIDAttr_InitsCodelabTracker() {
    const analytics = new CodelabAnalytics();
    analytics.setAttribute('gaid', 'UA-123');
    // Need to mock as we don't have window.ga.getAll()
    const mockGetAll = mockControl.createFunctionMock('getAll');
    const mockCreate = mockControl.createFunctionMock('create');
    mockGetAll().$returns([]).$anyTimes();
    // Creates 2 trackers (because of codelab gaid attribute).
    mockCreate().$times(2);

    window['ga'] = (...args) => {
      if (['create', 'getAll'].indexOf(args[0]) !== -1) {
        window['ga'][args[0]]();
      }
    };
    window['ga']['getAll'] = mockGetAll;
    window['ga']['create'] = mockCreate;

    mockControl.$replayAll();

    document.body.appendChild(analytics);
    analytics.setAttribute('codelab-gaid', 'UA-456');
    mockControl.$verifyAll();
  },

  async testSetAnalyticsReadyAttrs() {
    const analytics = new CodelabAnalytics();
    analytics.setAttribute('gaid', 'UA-123');
    // Need to mock as we don't have window.ga.getAll()
    const mockGetAll = mockControl.createFunctionMock('getAll');
    const mockCreate = mockControl.createFunctionMock('create');
    mockGetAll().$returns([]).$anyTimes();
    // Creates 2 trackers (because of codelab gaid attribute).
    mockCreate().$times(2);

    window['ga'] = (...args) => {
      if (['create', 'getAll'].indexOf(args[0]) !== -1) {
        window['ga'][args[0]]();
      }
    };
    window['ga']['getAll'] = mockGetAll;
    window['ga']['create'] = mockCreate;

    mockControl.$replayAll();

    const codelabElement = document.createElement('google-codelab');
    document.body.appendChild(codelabElement);
    document.body.appendChild(analytics);

    // This is obviously awful, but for some reason
    // mockControl.$waitAndVerifyAll() isn't working with closure_js_test. See
    // https://github.com/bazelbuild/rules_closure/issues/316. Once that's
    // resolved we can use it in place of the timeout.
    setTimeout(() => {
      assertEquals('', codelabElement.getAttribute('analytics-ready'));
    }, 5000);

  },

  testPageviewEventDispatch_SendsPageViewTracking() {

  },

  testEventDispatch_SendsEventTracking() {

  },

  testCodelabAttributes_UpdatesTrackingParams() {

  }
});
