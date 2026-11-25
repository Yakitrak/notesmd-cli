"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = isSpecifierDefault;
var _index = require("./generated/index.js");
function isSpecifierDefault(specifier) {
  if ((0, _index.isImportDefaultSpecifier)(specifier)) return true;
  const name = specifier.imported || specifier.exported;
  return (0, _index.isIdentifier)(name, {
    name: "default"
  }) || (0, _index.isStringLiteral)(name, {
    value: "default"
  });
}

//# sourceMappingURL=isSpecifierDefault.js.map
