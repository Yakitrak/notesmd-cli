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
export function toHexadecimal(code: number, next: number, omit: boolean | null | undefined): string;
//# sourceMappingURL=to-hexadecimal.d.ts.map