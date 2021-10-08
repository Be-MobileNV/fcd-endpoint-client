# openapi_client.DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**get_ws**](DefaultApi.md#get_ws) | **GET** /ws | websocket to push fcd data to


# **get_ws**
> get_ws(fcd=fcd)

websocket to push fcd data to

used to push fcd-data to endpoint

### Example

```python
from __future__ import print_function
import time
import openapi_client
from openapi_client.rest import ApiException
from pprint import pprint

# Create an instance of the API class
api_instance = openapi_client.DefaultApi()
fcd = openapi_client.Fcd() # Fcd |  (optional)

try:
    # websocket to push fcd data to
    api_instance.get_ws(fcd=fcd)
except ApiException as e:
    print("Exception when calling DefaultApi->get_ws: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **fcd** | [**Fcd**](Fcd.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

