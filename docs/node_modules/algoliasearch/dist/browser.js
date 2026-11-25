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

// builds/browser.ts
var browser_exports = {};
__export(browser_exports, {
  algoliasearch: () => algoliasearch,
  apiClientVersion: () => apiClientVersion
});
import { DEFAULT_REPLACE_ALL_OBJECTS_MAX_RETRIES, logWarning, withRequestId } from "@algolia/client-common";
import { abtestingV3Client } from "@algolia/abtesting";
import { abtestingClient } from "@algolia/client-abtesting";
import { analyticsClient } from "@algolia/client-analytics";
import { insightsClient } from "@algolia/client-insights";
import { personalizationClient } from "@algolia/client-personalization";
import { querySuggestionsClient } from "@algolia/client-query-suggestions";
import { searchClient } from "@algolia/client-search";
import { ingestionClient } from "@algolia/ingestion";
import { monitoringClient } from "@algolia/monitoring";
import { recommendClient } from "@algolia/recommend";

// builds/models.ts
var models_exports = {};
__export(models_exports, {
  apiClientVersion: () => apiClientVersion
});
__reExport(models_exports, abtesting_star);
__reExport(models_exports, client_abtesting_star);
__reExport(models_exports, client_analytics_star);
__reExport(models_exports, client_insights_star);
__reExport(models_exports, client_personalization_star);
__reExport(models_exports, client_query_suggestions_star);
__reExport(models_exports, client_search_star);
__reExport(models_exports, ingestion_star);
__reExport(models_exports, monitoring_star);
__reExport(models_exports, recommend_star);
import { apiClientVersion } from "@algolia/client-search";
import * as abtesting_star from "@algolia/abtesting";
import * as client_abtesting_star from "@algolia/client-abtesting";
import * as client_analytics_star from "@algolia/client-analytics";
import * as client_insights_star from "@algolia/client-insights";
import * as client_personalization_star from "@algolia/client-personalization";
import * as client_query_suggestions_star from "@algolia/client-query-suggestions";
import * as client_search_star from "@algolia/client-search";
import * as ingestion_star from "@algolia/ingestion";
import * as monitoring_star from "@algolia/monitoring";
import * as recommend_star from "@algolia/recommend";

// builds/browser.ts
__reExport(browser_exports, models_exports);
function algoliasearch(appId, apiKey, options) {
  if (!appId || typeof appId !== "string") {
    throw new Error("`appId` is missing.");
  }
  if (!apiKey || typeof apiKey !== "string") {
    throw new Error("`apiKey` is missing.");
  }
  const client = searchClient(appId, apiKey, options);
  let transformationConfig;
  if (options?.transformationOptions) {
    transformationConfig = options.transformationOptions;
  } else if (options?.transformation) {
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
    ingestionTransporter = ingestionClient(appId, apiKey, region, ingestionOptions);
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
      maxRetries = DEFAULT_REPLACE_ALL_OBJECTS_MAX_RETRIES
    }, requestOptions) {
      if (!ingestionTransporter) {
        throw new Error(
          "`transformationOptions` must be set in the client config before calling this method. It defaults to the Ingestion API defaults. See https://www.algolia.com/doc/libraries/sdk/methods/ingestion/"
        );
      }
      const searchRequestOptions = withRequestId(client.transporter, requestOptions);
      if (objects.length === 0) {
        logWarning(
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
      return abtestingClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initAbtestingV3: (initOptions) => {
      return abtestingV3Client(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initAnalytics: (initOptions) => {
      return analyticsClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initIngestion: (initOptions) => {
      return ingestionClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initInsights: (initOptions) => {
      return insightsClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initMonitoring: (initOptions = {}) => {
      return monitoringClient(initOptions.appId || appId, initOptions.apiKey || apiKey, initOptions.options);
    },
    initPersonalization: (initOptions) => {
      return personalizationClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initQuerySuggestions: (initOptions) => {
      return querySuggestionsClient(
        initOptions.appId || appId,
        initOptions.apiKey || apiKey,
        initOptions.region,
        initOptions.options
      );
    },
    initRecommend: (initOptions = {}) => {
      return recommendClient(initOptions.appId || appId, initOptions.apiKey || apiKey, initOptions.options);
    }
  };
}
export {
  algoliasearch,
  apiClientVersion
};
//# sourceMappingURL=browser.js.map