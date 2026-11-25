/**
 * @import {FormatSmartOptions} from './util/format-smart.js'
 */

/**
 * @typedef LightOptions
 *   Configuration for `stringifyEntitiesLight`.
 * @property {boolean | null | undefined} [escapeOnly=false]
 *   Whether to only escape possibly dangerous characters (default: `false`);
 *   those characters are `"`, `&`, `'`, `<`, `>`, and `` ` ``.
 * @property {ReadonlyArray<string> | null | undefined} [subset]
 *   Whether to only escape the given subset of characters (optional);
 *   only BMP characters are supported here (so no emoji).
 */

/**
 * @typedef FormatOptions
 *   Internal format options.
 * @property {(code: number, next: number, options: CoreWithFormatOptions) => string} format
 *   Format strategy.
 */

/**
 * @typedef {LightOptions & FormatOptions & FormatSmartOptions} CoreWithFormatOptions
 */

const defaultSubsetRegex = /["&'<>`]/g
const surrogatePairsRegex = /[\uD800-\uDBFF][\uDC00-\uDFFF]/g
const controlCharactersRegex =
  // eslint-disable-next-line no-control-regex, regexp/hexadecimal-escape, regexp/no-control-character, regexp/no-obscure-range -- matching control characters is the point.
  /[\x01-\t\v\f\x0E-\x1F\x7F\x81\x8D\x8F\x90\x9D\xA0-\uFFFF]/g
const regexEscapeRegex = /[|\\{}()[\]^$+*?.]/g

/** @type {WeakMap<ReadonlyArray<string>, RegExp>} */
const subsetToRegexCache = new WeakMap()

/**
 * Encode certain characters in `value`.
 *
 * @param {string} value
 *   Value to encode.
 * @param {CoreWithFormatOptions} options
 *   Configuration.
 * @returns {string}
 *   Encoded value.
 */
export function core(value, options) {
  value = value.replace(
    options.subset
      ? charactersToExpressionCached(options.subset)
      : defaultSubsetRegex,
    basic
  )

  if (options.subset || options.escapeOnly) {
    return value
  }

  return (
    value
      // Surrogate pairs.
      .replace(surrogatePairsRegex, surrogate)
      // BMP control characters (C0 except for LF, CR, SP; DEL; and some more
      // non-ASCII ones).
      .replace(controlCharactersRegex, basic)
  )

  /**
   * Encode a surrogate pair.
   *
   * @param {string} pair
   *   Surrogate pair.
   * @param {number} index
   *   Index of `pair` in `all`.
   * @param {string} all
   *   Whole value.
   * @returns {string}
   *   Character reference.
   */
  function surrogate(pair, index, all) {
    return options.format(
      (pair.charCodeAt(0) - 0xd800) * 0x400 +
        pair.charCodeAt(1) -
        0xdc00 +
        0x10000,
      all.charCodeAt(index + 2),
      options
    )
  }

  /**
   * Encode a character.
   *
   * @param {string} character
   *   Character.
   * @param {number} index
   *   Index of `character` in `all`.
   * @param {string} all
   *   Whole value.
   * @returns {string}
   *   Character reference.
   */
  function basic(character, index, all) {
    return options.format(
      character.charCodeAt(0),
      all.charCodeAt(index + 1),
      options
    )
  }
}

/**
 * A wrapper function that caches the result of `charactersToExpression` with a WeakMap.
 * This can improve performance when tooling calls `charactersToExpression` repeatedly
 * with the same subset.
 *
 * @param {ReadonlyArray<string>} subset
 *   Characters to match.
 * @returns {RegExp}
 *   Expression.
 */
function charactersToExpressionCached(subset) {
  let cached = subsetToRegexCache.get(subset)

  if (!cached) {
    cached = charactersToExpression(subset)
    subsetToRegexCache.set(subset, cached)
  }

  return cached
}

/**
 * Create an expression that matches any of `subset`.
 *
 * @param {ReadonlyArray<string>} subset
 *   Characters to match.
 * @returns {RegExp}
 *   Expression.
 */
function charactersToExpression(subset) {
  /** @type {Array<string>} */
  const groups = []
  let index = -1

  while (++index < subset.length) {
    groups.push(subset[index].replace(regexEscapeRegex, '\\$&'))
  }

  return new RegExp('(?:' + groups.join('|') + ')', 'g')
}
