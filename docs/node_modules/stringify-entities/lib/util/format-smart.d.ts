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
export function formatSmart(code: number, next: number, options: FormatSmartOptions): string;
/**
 * Configuration for `stringifyEntities` that is not in `LightOptions`.
 */
export type FormatSmartOptions = {
    /**
     * Create character references which don’t fail in attributes
     * (default: `false`);
     * only applies when operating dangerously with
     * `omitOptionalSemicolons: true`.
     */
    attribute?: boolean | null | undefined;
    /**
     * Whether to omit semicolons when possible (default: `false`);
     * this creates what HTML calls “parse errors” but is otherwise still valid
     * HTML,
     * so don’t use this except when building a minifier;
     * omitting semicolons is possible for certain named and numeric references in
     * some cases.
     */
    omitOptionalSemicolons?: boolean | null | undefined;
    /**
     * Prefer named character references (`&amp;`) where possible
     * (default: `false`).
     */
    useNamedReferences?: boolean | null | undefined;
    /**
     * Prefer the shortest possible reference,
     * if that results in less bytes (default: `false`);
     * `useNamedReferences` can be omitted when using `useShortestReferences`.
     */
    useShortestReferences?: boolean | null | undefined;
};
//# sourceMappingURL=format-smart.d.ts.map