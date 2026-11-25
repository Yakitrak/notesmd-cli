/**
 * @typedef FormatSmartOptions
 *   Configuration for `stringifyEntities` that is not in `LightOptions`.
 * @property {boolean | null | undefined} [attribute=false]
 *   Create character references which don’t fail in attributes
 *   (default: `false`);
 *   only applies when operating dangerously with
 *   `omitOptionalSemicolons: true`.
 * @property {boolean | null | undefined} [omitOptionalSemicolons=false]
 *   Whether to omit semicolons when possible (default: `false`);
 *   this creates what HTML calls “parse errors” but is otherwise still valid
 *   HTML,
 *   so don’t use this except when building a minifier;
 *   omitting semicolons is possible for certain named and numeric references in
 *   some cases.
 * @property {boolean | null | undefined} [useNamedReferences=false]
 *   Prefer named character references (`&amp;`) where possible
 *   (default: `false`).
 * @property {boolean | null | undefined} [useShortestReferences=false]
 *   Prefer the shortest possible reference,
 *   if that results in less bytes (default: `false`);
 *   `useNamedReferences` can be omitted when using `useShortestReferences`.
 */

import {toHexadecimal} from './to-hexadecimal.js'
import {toDecimal} from './to-decimal.js'
import {toNamed} from './to-named.js'

/**
 * Configurable ways to encode a character yielding pretty or small results.
 *
 * @param {number} code
 *   Code point.
 * @param {number} next
 *   Code point of the next character (or `NaN`).
 * @param {FormatSmartOptions} options
 *   Configuration.
 * @returns {string}
 *   Character reference.
 */
export function formatSmart(code, next, options) {
  const omit = options.omitOptionalSemicolons && next !== 59 /* `;` */
  let numeric = toHexadecimal(code, next, omit)
  /** @type {string | undefined} */
  let named

  if (options.useNamedReferences || options.useShortestReferences) {
    named = toNamed(code, next, omit, options.attribute)
  }

  // Use the shortest numeric reference when requested.
  // A simple algorithm would use decimal for all code points under 100, as
  // those are shorter than hexadecimal:
  //
  // * `&#99;` vs `&#x63;` (decimal shorter)
  // * `&#100;` vs `&#x64;` (equal)
  //
  // However, because we take `next` into consideration when `omit` is used,
  // And it would be possible that decimals are shorter on bigger values as
  // well if `next` is hexadecimal but not decimal, we instead compare both.
  if (
    (options.useShortestReferences || !named) &&
    options.useShortestReferences
  ) {
    const decimal = toDecimal(code, next, omit)

    if (decimal.length < numeric.length) {
      numeric = decimal
    }
  }

  return named &&
    (!options.useShortestReferences || named.length < numeric.length)
    ? named
    : numeric
}
