goog.module('googlecodelabs.CodelabAnalyticsDef');
const CodelabAnalytics = goog.require('googlecodelabs.CodelabAnalytics');

try {
  window.customElements.define(CodelabAnalytics.getTagName(), CodelabAnalytics);
} catch (e) {
  console.warn('googlecodelabs.CodelabAnalytics', e);
}