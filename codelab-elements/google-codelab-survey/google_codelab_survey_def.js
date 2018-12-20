goog.module('googlecodelabs.CodelabSurveyDef');
const CodelabSurvey = goog.require('googlecodelabs.CodelabSurvey');

try {
  window.customElements.define(CodelabSurvey.getTagName(), CodelabSurvey);
} catch (e) {
  console.warn('googlecodelabs.CodelabSurvey', e);
}