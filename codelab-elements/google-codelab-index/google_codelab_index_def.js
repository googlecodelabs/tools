goog.module('googlecodelabs.CodelabIndexDef');
const CodelabIndex = goog.require('googlecodelabs.CodelabIndex');

try {
  window.customElements.define(CodelabIndex.getTagName(), CodelabIndex);
} catch (e) {
  console.warn('googlecodelabs.CodelabIndex', e);
}