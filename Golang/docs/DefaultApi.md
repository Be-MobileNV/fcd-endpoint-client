# \DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetWs**](DefaultApi.md#GetWs) | **Get** /ws | Your GET endpoint



## GetWs

> GetWs(ctx, optional)

Your GET endpoint

used to push gps-positions to endpoint

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetWsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a GetWsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **gpsPosition** | [**optional.Interface of GpsPosition**](GpsPosition.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

