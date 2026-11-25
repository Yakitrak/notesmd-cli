const decimalRegex = /\d/

/**
 * Configurable ways to encode characters as decimal references.
 *
 * @param {number} code
 *   Code point.
 * @param {number} next
 *   Code point of the next character (or `NaN`).
 * @param {boolean | null | undefined} omit
 *   Whether to omit the semicolon when possible.
 * @returns {string}
 *   Decimal reference.
 */
export function toDecimal(code, next, omit) {
  const value = '&#' + String(code)
  return omit && next && !decimalRegex.test(String.fromCharCode(next))
    ? value
    : value + ';'
}
