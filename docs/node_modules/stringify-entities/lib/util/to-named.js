import {characterEntitiesLegacy} from 'character-entities-legacy'
import {characterEntitiesHtml4} from 'character-entities-html4'
import {dangerous} from '../constant/dangerous.js'

const own = {}.hasOwnProperty

/**
 * `characterEntitiesHtml4` but inverted.
 *
 * @type {Record<string, string>}
 */
const characters = {}

/** @type {string} */
let key

for (key in characterEntitiesHtml4) {
  if (own.call(characterEntitiesHtml4, key)) {
    characters[characterEntitiesHtml4[key]] = key
  }
}

const notAlphanumericRegex = /[^\dA-Za-z]/

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
export function toNamed(code, next, omit, attribute) {
  const character = String.fromCharCode(code)

  if (own.call(characters, character)) {
    const name = characters[character]
    const value = '&' + name

    if (
      omit &&
      characterEntitiesLegacy.includes(name) &&
      !dangerous.includes(name) &&
      (!attribute ||
        (next &&
          next !== 61 /* `=` */ &&
          notAlphanumericRegex.test(String.fromCharCode(next))))
    ) {
      return value
    }

    return value + ';'
  }

  return ''
}
