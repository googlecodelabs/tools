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

goog.module('googlecodelabs.CodelabCodeMenu');

const Templates = goog.require('googlecodelabs.CodelabCodeMenu.Templates');
const soy = goog.require('goog.soy');


/**
 * @extends {HTMLElement}
 */
class CodelabCodeMenu extends HTMLElement {
  /** @return {string} */
  static getTagName() { return 'google-codelab-code-menu'; }


  constructor() {
    super();

    this.addEventListener('click', e => {

      const inputEl = /** @type {?HTMLInputElement} */ (e.target);
      switch (inputEl.title) {
        case "Wrap":
          this.wrap_(inputEl);

          break;

        case "Theme":
          this.theme_(inputEl);

          break;

        case "Copy":
          this.copy_(inputEl);

          break;
      }
    });
  }


  /**
   * @param {!HTMLInputElement} el 
   * @private 
   */
  wrap_(el) {
    //console.log( el.getTagName());
    const preEl = /** @type {?Element} */ (el.parentNode.parentNode.parentNode);
    const style = /** @type {?CSSStyleDeclaration} */ (preEl.getElementsByTagName("code")[0].style);

    if (style.whiteSpace == "pre-wrap") {
      style.removeProperty("white-space");
      // wrap
      el.value = "\u21E5";
    } else {
      style.whiteSpace = "pre-wrap";
      // unwrap
      el.value = "\u2194";
    }
  }

  /**
   * Convert rgb[a](d, d, d[, d]) color into #hhhhhh[hh] notation
   * 
   * @param {string} x
   * @return {string} 
   * @private
  */
  h_(x) { return '#' + x.match(/\d+/g).map((z) => ((+z < 16) ? '0' : '') + (+z).toString(16)).join(''); }


  /**
   * @param {!HTMLInputElement} el 
   * @private 
   */
  theme_(el) {
    const preEl = /** @type {?Element} */ (el.parentNode.parentNode.parentNode);
    const preStyle = /** @type {?CSSStyleDeclaration} */ (getComputedStyle(preEl));

    const ctheme = this.h_(preStyle.backgroundColor) == "#28323f" ? "Light" : "Dark";

    for (let i = 0; i < document.styleSheets.length; i++) {
      if (document.styleSheets[i].href && document.styleSheets[i].href.endsWith("/codelab-elements.css")) {
        const ss = /** @type {?CSSStyleSheet} */ (document.styleSheets[i]);

        // google-codelab-step pre
        //  background-color: #185abc;
        // 
        for (let r = 0; r < ss.cssRules.length; r++) {
          const rule = /** @type {!CSSStyleRule} */ (ss.cssRules[r]);

          switch (rule.selectorText) {
            case "google-codelab-step pre":
              rule.style.backgroundColor = ctheme == "Dark" ? "#28323f" : "#f8f9fa";
              rule.style.color = ctheme == "Dark" ? "#ffffff" : "#000000";
              break;
            case "google-codelab-step pre .pln, google-codelab-step code .pln":
              rule.style.color = ctheme == "Dark" ? "#F8F9FA" : "#000000";
              break;

            case "google-codelab-step pre .pun, google-codelab-step code .pun":
              rule.style.color = ctheme == "Dark" ? "#F8F9FA" : "#1e8e3e";
              break;
          }
        }
      }
    }
  }

  /**
   * @param {!HTMLInputElement} el 
   * @private 
   */
  copy_(el) {
    var copyEl = el;
    const preEl = /** @type {?Element} */ (el.parentNode.parentNode.parentNode);
    const text = preEl.children[1].textContent;
    navigator.clipboard.writeText(text)
      .then(() => {
        console.log('Text copied to clipboard');
        copyEl.value = "\u220E";
        setTimeout(() => {
          copyEl.value = "\u29c9";
        }, 250);
      })
      .catch(err => {
        // This can happen if the user denies clipboard permissions:
        console.error('Could not copy text: ', err);
      });
  }

  /**
   * @export
   * @override
   */
  connectedCallback() {
    this.updateDom_();
  }

  /** @private */
  updateDom_() {
    const updatedDom = soy.renderAsElement(Templates.codeMenu, {});
    this.appendChild(updatedDom);
  }
}

exports = CodelabCodeMenu;
