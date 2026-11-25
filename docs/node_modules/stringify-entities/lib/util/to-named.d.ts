/**
 * Configurable ways to encode characters as named references.
 *
 * @param {number} code
 *   Code point.
 * @param {number} next
 *   Code point of the next character (or `NaN`).
 * @param {boolean | null | undefined} omit
 *   Whether to omit the semicolon when possible.
 * @param {boolean | null | undefined} attribute
 *   Whether the reference is in an attribute.
 * @returns {string}
 *   Named reference, or empty string if there is no name.
 */
export function toNamed(code: number, next: number, omit: boolean | null | undefined, attribute: boolean | null | undefined): string;
//# sourceMappingURL=to-named.d.ts.map