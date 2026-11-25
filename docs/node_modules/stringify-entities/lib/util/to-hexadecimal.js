const hexadecimalRegex = /[\dA-Fa-f]/

/**
 * Configurable ways to encode characters as hexadecimal references.
 *
 * @param {number} code
 *   Code point.
 * @param {number} next
 *   Code point of the next character (or `NaN`).
 * @param {boolean | null | undefined} omit
 *   Whether to omit the semicolon when possible.
 * @returns {string}
 *   Hexadecimal reference.
 */
export function toHexadecimal(code, next, omit) {
  const value = '&#x' + code.toString(16).toUpperCase()
  return omit && next && !hexadecimalRegex.test(String.fromCharCode(next))
    ? value
    : value + ';'
}
