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
export function toDecimal(code: number, next: number, omit: boolean | null | undefined): string;
//# sourceMappingURL=to-decimal.d.ts.map