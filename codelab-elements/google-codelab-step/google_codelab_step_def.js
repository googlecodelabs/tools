goog.module('googlecodelabs.CodelabStepDef');
const CodelabStep = goog.require('googlecodelabs.CodelabStep');

try {
  window.customElements.define(CodelabStep.getTagName(), CodelabStep);
} catch (e) {
  console.warn('googlecodelabs.CodelabStep', e);
}