goog.module('googlecodelabs.CodelabIndex.CardsDef');
const CodelabIndexCards = goog.require('googlecodelabs.CodelabIndex.Cards');

try {
  window.customElements.define(CodelabIndexCards.getTagName(), CodelabIndexCards);
} catch (e) {
  console.warn('googlecodelabs.CodelabIndex.Cards', e);
}