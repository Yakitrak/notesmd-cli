// src/createHttpRequester.ts
import http from "http";
import https from "https";
import { Readable } from "stream";
import { URL } from "url";
import zlib from "zlib";
import { StreamRequestError } from "@algolia/client-common";
var agentOptions = { keepAlive: true };
var defaultHttpAgent = new http.Agent(agentOptions);
var defaultHttpsAgent = new https.Agent(agentOptions);
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
      const url = new URL(request.url);
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
      const req = (url.protocol === "https:" ? https : http).request(options, (response) => {
        let contentBuffers = [];
        response.on("data", (chunk) => {
          contentBuffers = contentBuffers.concat(chunk);
        });
        response.on("end", () => {
          clearTimeout(connectTimeout);
          clearTimeout(responseTimeout);
          let buffer = Buffer.concat(contentBuffers);
          if (response.headers["content-encoding"] === "gzip") {
            buffer = zlib.gunzipSync(buffer);
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
      const url = new URL(request.url);
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
      const req = (url.protocol === "https:" ? https : http).request(options, (response) => {
        const statusCode = response.statusCode || 0;
        if (statusCode < 200 || statusCode >= 300) {
          let body = "";
          response.on("data", (chunk) => {
            body += chunk;
          });
          response.on("end", () => {
            reject(new StreamRequestError(statusCode, body, toResponseHeaders(response.headers)));
          });
          return;
        }
        resolve(Readable.toWeb(response));
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
export {
  createHttpRequester
};
//# sourceMappingURL=requester.http.js.map