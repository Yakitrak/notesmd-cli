/**
 * The smallest way to encode a character.
 *
 * @param {number} code
 *   Code point.
 * @returns {string}
 *   Hexadecimal reference.
 */
export function formatBasic(code) {
  return '&#x' + code.toString(16).toUpperCase() + ';'
}
