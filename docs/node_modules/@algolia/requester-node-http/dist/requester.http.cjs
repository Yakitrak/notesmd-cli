"use strict";
var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var src_exports = {};
__export(src_exports, {
  createHttpRequester: () => createHttpRequester
});
module.exports = __toCommonJS(src_exports);

// src/createHttpRequester.ts
var import_http = __toESM(require("http"), 1);
var import_https = __toESM(require("https"), 1);
var import_stream = require("stream");
var import_url = require("url");
var import_zlib = __toESM(require("zlib"), 1);
var import_client_common = require("@algolia/client-common");
var agentOptions = { keepAlive: true };
var defaultHttpAgent = new import_http.default.Agent(agentOptions);
var defaultHttpsAgent = new import_https.default.Agent(agentOptions);
function toResponseHeaders(incomingHeaders) {
  const headers = {};
  for (const [name, value] of Object.entries(incomingHeaders)) {
    if (value !== void 0) {
      headers[name] = Array.isArray(value) ? value.join(", ") : value;
    }
  }
  return headers;
}
function createHttpRequester({
  agent: userGlobalAgent,
  httpAgent: userHttpAgent,
  httpsAgent: userHttpsAgent,
  requesterOptions = {}
} = {}) {
  const httpAgent = userHttpAgent || userGlobalAgent || defaultHttpAgent;
  const httpsAgent = userHttpsAgent || userGlobalAgent || defaultHttpsAgent;
  function send(request) {
    return new Promise((resolve) => {
      let responseTimeout;
      let connectTimeout;
      const url = new import_url.URL(request.url);
      const path = url.search === null ? url.pathname : `${url.pathname}${url.search}`;
      const privateHeaders = {
        "accept-encoding": "gzip"
      };
      if (request.data !== void 0 && request.method === "DELETE") {
        privateHeaders["content-length"] = String(
          typeof request.data === "string" ? Buffer.byteLength(request.data) : request.data.byteLength
        );
      }
      const options = {
        agent: url.protocol === "https:" ? httpsAgent : httpAgent,
        hostname: url.hostname,
        path,
        method: request.method,
        ...requesterOptions,
        headers: {
          ...privateHeaders,
          ...request.headers,
          ...requesterOptions.headers
        }
      };
      if (url.port && !requesterOptions.port) {
        options.port = url.port;
      }
      const req = (url.protocol === "https:" ? import_https.default : import_http.default).request(options, (response) => {
        let contentBuffers = [];
        response.on("data", (chunk) => {
          contentBuffers = contentBuffers.concat(chunk);
        });
        response.on("end", () => {
          clearTimeout(connectTimeout);
          clearTimeout(responseTimeout);
          let buffer = Buffer.concat(contentBuffers);
          if (response.headers["content-encoding"] === "gzip") {
            buffer = import_zlib.default.gunzipSync(buffer);
          }
          resolve({
            status: response.statusCode || 0,
            content: buffer.toString(),
            headers: toResponseHeaders(response.headers),
            isTimedOut: false
          });
        });
        response.on("error", (error) => {
          clearTimeout(connectTimeout);
          clearTimeout(responseTimeout);
          resolve({ status: 0, content: error.message, isTimedOut: false });
        });
      });
      const createTimeout = (timeout, content) => {
        return setTimeout(() => {
          req.destroy();
          resolve({
            status: 0,
            content,
            isTimedOut: true
          });
        }, timeout);
      };
      connectTimeout = createTimeout(request.connectTimeout, "Connection timeout");
      req.on("error", (error) => {
        clearTimeout(connectTimeout);
        clearTimeout(responseTimeout);
        resolve({ status: 0, content: error.message, isTimedOut: false });
      });
      req.once("response", () => {
        clearTimeout(connectTimeout);
        responseTimeout = createTimeout(request.responseTimeout, "Socket timeout");
      });
      if (request.data !== void 0) {
        req.write(request.data);
      }
      req.end();
    });
  }
  function sendStream(request) {
    return new Promise((resolve, reject) => {
      const url = new import_url.URL(request.url);
      const path = url.search === null ? url.pathname : `${url.pathname}${url.search}`;
      const options = {
        agent: url.protocol === "https:" ? httpsAgent : httpAgent,
        hostname: url.hostname,
        path,
        method: request.method,
        ...requesterOptions,
        headers: {
          ...request.headers,
          ...requesterOptions.headers
        }
      };
      if (url.port && !requesterOptions.port) {
        options.port = url.port;
      }
      const req = (url.protocol === "https:" ? import_https.default : import_http.default).request(options, (response) => {
        const statusCode = response.statusCode || 0;
        if (statusCode < 200 || statusCode >= 300) {
          let body = "";
          response.on("data", (chunk) => {
            body += chunk;
          });
          response.on("end", () => {
            reject(new import_client_common.StreamRequestError(statusCode, body, toResponseHeaders(response.headers)));
          });
          return;
        }
        resolve(import_stream.Readable.toWeb(response));
      });
      req.on("error", (error) => {
        reject(error);
      });
      if (request.data !== void 0) {
        req.write(request.data);
      }
      req.end();
    });
  }
  return { send, sendStream };
}
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  createHttpRequester
});
//# sourceMappingURL=requester.http.cjs.map