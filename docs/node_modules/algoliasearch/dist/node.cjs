"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
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
var __reExport = (target, mod, secondTarget) => (__copyProps(target, mod, "default"), secondTarget && __copyProps(secondTarget, mod, "default"));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// builds/node.ts
var node_exports = {};
__export(node_exports, {
  algoliasearch: () => algoliasearch,
  apiClientVersion: () => import_client_search.apiClientVersion
});
module.exports = __toCommonJS(node_exports);
var import_client_common = require("@algolia/client-common");
var import_abtesting = require("@algolia/abtesting");
var import_client_abtesting = require("@algolia/client-abtesting");
var import_client_analytics = require("@algolia/client-analytics");
var import_client_insights = require("@algolia/client-insights");
var import_client_personalization = require("@algolia/client-personalization");
var import_client_query_suggestions = require("@algolia/client-query-suggestions");
var import_client_search2 = require("@algolia/client-search");
var import_ingestion = require("@algolia/ingestion");
var import_monitoring = require("@algolia/monitoring");
var import_recommend = require("@algolia/recommend");

// builds/models.ts
var models_exports = {};
__export(models_exports, {
  apiClientVersion: () => import_client_search.apiClientVersion
});
var import_client_search = require("@algolia/client-search");
__reExport(models_exports, require("@algolia/abtesting"));
__reExport(models_exports, require("@algolia/client-abtesting"));
__reExport(models_exports, require("@algolia/client-analytics"));
__reExport(models_exports, require("@algolia/client-insights"));
__reExport(models_exports, require("@algolia/client-personalization"));
__reExport(models_exports, require("@algolia/client-query-suggestions"));
__reExport(models_exports, require("@algolia/client-search"));
__reExport(models_exports, require("@algolia/ingestion"));
__reExport(models_exports, require("@algolia/monitoring"));
__reExport(models_exports, require("@algolia/recommend"));

// builds/node.ts
__reExport(node_exports, models_exports, module.exports);
function algoliasearch(appId, apiKey, options) {
  if (!appId || typeof appId !== "string") {
    throw new Error("`appId` is missing.");
  }
  if (!apiKey || typeof apiKey !== "string") {
    throw new Error("`apiKey` is missing.");
  }
  const client = (0, import_client_search2.searchClient)(appId, apiKey, options);
  let transformationConfig;
  if (options == null ? void 0 : options.transformationOptions) {
    transformationConfig = options.transformationOptions;
  } else if (options == null ? void 0 : options.transformation) {
    transformationConfig = { region: options.transformation.region };
  }
  let ingestionTransporter;
  if (transformationConfig) {
    if (!transformationConfig.region) {
      throw new Error(
        "`region` is required in `transformationOptions`. See https://www.algolia.com/doc/libraries/sdk/methods/ingestion/"
      );
    }
    const { region, ...ingestionOptions } = transformationConfig;
    ingestionTransporter = (0, import_ingestion.ingestionClient)(appId, apiKey, region, ingestionOptions);
  }
  return {
    ...client,
    async saveObjectsWithTransformation({ indexName, objects, batchSize, waitForTasks, maxRetries }, requestOptions) {
      if (!ingestionTransporter) {
        throw new Error(
          "`transformationOptions` must be set in the client config before calling this method. It defaults to the Ingestion API defaults. See https://www.algolia.com/doc/libraries/sdk/methods/ingestion/"
        );
      }
      return ingestionTransporter.chunkedPush(
        { indexName, objects, action: "addObject", batchSize, waitForTasks, maxRetries },
        requestOptions
      );
    },
    async partialUpdateObjectsWithTransformation({ indexName, objects, createIfNotExists, batchSize, waitForTasks, maxRetries }, requestOptions) {
      if (!ingestionTransporter) {
        throw new Error(
          "`transformationOptions` must be set in the client config before calling this method. It defaults to the Ingestion API defaults. See https://www.algolia.com/doc/libraries/sdk/methods/ingestion/"
        );
      }
      return ingestionTransporter.chunkedPush(
        {
          indexName,
          objects,
          action: createIfNotExists ? "partialUpdateObject" : "partialUpdateObjectNoCreate",
          batchSize,
          waitForTasks,
          maxRetries
        },
        requestOptions
      );
    },
    async replaceAllObjectsWithTransformation({
      indexName,
      objects,
      batchSize,
      scopes,
      maxRetries = import_client_common.DEFAULT_REPLACE_ALL_OBJECTS_MAX_RETRIES
    }, requestOptions) {
      if (!ingestionTransporter) {
        throw new Error(
          "`transformationOptions` must be set in the client config before calling this method. It defaults to the Ingestion API defaults. See https://www.algolia.com/doc/libraries/sdk/methods/ingestion/"
        );
      }
      const searchRequestOptions = (0, import_client_common.withRequestId)(client.transporter, requestOptions);
      if (objects.length === 0) {
        (0, import_client_common.logWarning)(
          client.transporter.logger,
          `replaceAllObjectsWithTransformation was called with an empty list of objects, which will delete all records currently in the "${indexName}" index.`
        );
      }
      const randomSuffix = Math.floor(Math.random() * 1e6) + 1e5;
      const tmpIndexName = `${indexName}_tmp_${randomSuffix}`;
      if (scopes === void 0) {
        scopes = ["settings", "rules", "synonyms"];
      }
      try {
        let copyOperationResponse = await this.operationIndex(
          {
            indexName,
            operationIndexParams: {
              operation: "copy",
              destination: tmpIndexName,
              scope: scopes
            }
          },
          searchRequestOptions
        );
        const watchResponses = await ingestionTransporter.chunkedPush(
          {
            indexName: tmpIndexName,
            objects,
            waitForTasks: true,
            batchSize,
            referenceIndexName: indexName,
            maxRetries
          },
          requestOptions
        );
        await this.waitForTask(
          {
            indexName: tmpIndexName,
            taskID: copyOperationResponse.taskID,
            maxRetries
          },
          searchRequestOptions
        );
        copyOperationResponse = await this.operationIndex(
          {
            indexName,
            operationIndexParams: {
              operation: "copy",
              destination: tmpIndexName,
              scope: scopes
            }
          },
          searchRequestOptions
        );
        await this.waitForTask(
          {
            indexName: tmpIndexName,
            taskID: copyOperationResponse.taskID,
            maxRetries
          },
          searchRequestOptions
        );
        const moveOperationResponse = await this.operationIndex(
          {
            indexName: tmpIndexName,
            operationIndexParams: { operation: "move", destination: indexName }
          },
          searchRequestOptions
        );
        await this.waitForTask(
          {
            indexName: tmpIndexName,
            taskID: moveOperationResponse.taskID,
            maxRetries
          },
          searchRequestOptions
        );
        return { copyOperationResponse, watchResponses, moveOperationResponse };
      } catch (error) {
        await this.deleteIndex({ indexName: tmpIndexName }, searchRequestOptions);
        throw error;
      }
    },
    /**
     * Get the value of the `algoliaAgent`, used by our libraries internally and telemetry system.
     */
    get _ua() {
      return client.transporter.algoliaAgent.value;
    },
    initAbtesting: (initOptions) => {
      return (0, import_client_abtesting.abtestingClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initAbtestingV3: (initOptions) => {
      return (0, import_abtesting.abtestingV3Client)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initAnalytics: (initOptions) => {
      return (0, import_client_analytics.analyticsClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initIngestion: (initOptions) => {
      return (0, import_ingestion.ingestionClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initInsights: (initOptions) => {
      return (0, import_client_insights.insightsClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initMonitoring: (initOptions = {}) => {
      return (0, import_monitoring.monitoringClient)(initOptions.appId || appId, initOptions.apiKey || apiKey, initOptions.options);
    },
    initPersonalization: (initOptions) => {
      return (0, import_client_personalization.personalizationClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initQuerySuggestions: (initOptions) => {
      return (0, import_client_query_suggestions.querySuggestionsClient)(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initRecommend: (initOptions = {}) => {
      return (0, import_recommend.recommendClient)(initOptions.appId || appId, initOptions.apiKey || apiKey, initOptions.options);
    }
  };
}
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  algoliasearch,
  apiClientVersion
});
//# sourceMappingURL=node.cjs.map