/* 
 * fcd-endpoint-client
 *
 * FCD-endpoint-client
 *
 * The version of the OpenAPI document: 1.0.0
 * Contact: api-support@be-mobile.com
 * Generated by: https://github.com/openapitools/openapi-generator.git
 */


using System;
using RestSharp;

namespace Org.OpenAPITools.Client
{
    /// <summary>
    /// A delegate to ExceptionFactory method
    /// </summary>
    /// <param name="methodName">Method name</param>
    /// <param name="response">Response</param>
    /// <returns>Exceptions</returns>
    public delegate Exception ExceptionFactory(string methodName, IRestResponse response);
}
