/**
 * @import {Parents, Text} from 'hast'
 * @import {Raw} from 'mdast-util-to-hast'
 * @import {State} from '../index.js'
 */

import {stringifyEntities} from 'stringify-entities'

// Declare array as variable so it can be cached by `stringifyEntities`
const textEntitySubset = ['<', '&']

// HTML elements whose text is not escaped.
// See: <https://html.spec.whatwg.org/multipage/parsing.html#serialising-html-fragments>.
const raw = new Set([
  'iframe',
  'noembed',
  'noframes',
  'plaintext',
  'script',
  'style',
  'xmp'
])

/**
 * Serialize a text node.
 *
 * @param {Raw | Text} node
 *   Node to handle.
 * @param {number | undefined} _
 *   Index of `node` in `parent`.
 * @param {Parents | undefined} parent
 *   Parent of `node`.
 * @param {State} state
 *   Info passed around about the current state.
 * @returns {string}
 *   Serialized node.
 */
export function text(node, _, parent, state) {
  // The schema is that of `parent`, as raw elements do not switch spaces.
  return parent &&
    parent.type === 'element' &&
    state.schema.space === 'html' &&
    raw.has(parent.tagName)
    ? node.value
    : stringifyEntities(
        node.value,
        Object.assign({}, state.settings.characterReferences, {
          subset: textEntitySubset
        })
      )
}
