package well_known

import (
	"encoding/json"
	"net/http"
	"os"
)

func HandleWellKnownEndpoint(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "issuer":                                           GetIssuer(),
        "authorization_endpoint":                           getAuthorizationEndpoint(),
        "token_endpoint":                                   getTokenEndpoint(),
        "jwks_uri":                                         getJwksUri(),
        "scopes_supported":                                 getScopesSupported(),
        "response_types_supported":                         getResponseTypesSupported(),
        "grant_types_supported":                            getGrantTypesSupported(),
        "token_endpoint_auth_methods_supported":            getTokenEndpointAuthMethodsSupported(),
        "code_challenge_methods_supported":                 getCodeChallengeMethodsSupported(),
        "token_endpoint_auth_signing_alg_values_supported": getTokenEndpointAuthSigningAlgValuesSupported(),
        "id_token_signing_alg_values_supported":            getIdTokenSigningAlgValuesSupported(),
        "request_object_signing_alg_values_supported":      getRequestObjectSigningAlgValuesSupported(),
        "response_modes_supported":                         getResponseModesSupported(),
        "subject_types_supported":                          getSubjectTypesSupported(),
        "userinfo_signing_alg_values_supported":            getUserinfoSigningAlgValuesSupported(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func GetIssuer() string {
    return getBaseUrl()
}

func getAuthorizationEndpoint() string {
    return getBaseUrl() + "/authorize"
}

func getTokenEndpoint() string {
    return getBaseUrl() + "/token"
}

func getJwksUri() string {
    return getBaseUrl() + "/jwks"
}

func getScopesSupported() []string {
    return []string{"read", "write", "delete"}
}

func getResponseTypesSupported() []string {
    return []string{"code", "token"}
}

func getGrantTypesSupported() []string {
    return []string{"authorization_code", "client_credentials", "refresh_token"}
}

func getTokenEndpointAuthMethodsSupported() []string {
    return []string{"client_secret_basic", "client_secret_post"}
}

func getCodeChallengeMethodsSupported() []string {
    return []string{"plain", "S256"}
}

func getTokenEndpointAuthSigningAlgValuesSupported() []string {
    return []string{"RS256"}
}

func getIdTokenSigningAlgValuesSupported() []string {
    return []string{"RS256"}
}

func getRequestObjectSigningAlgValuesSupported() []string {
    return []string{"RS256"}
}

func getResponseModesSupported() []string {
    return []string{"query", "fragment", "form_post"}
}

func getSubjectTypesSupported() []string {
    return []string{"public"}
}

func getUserinfoSigningAlgValuesSupported() []string {
    return []string{"RS256"}
}

func getBaseUrl() string {
    extUrl, exists := os.LookupEnv("EXTERNAL_URL")
    if exists {
        return extUrl
    }
    port,_ := os.LookupEnv("PORT")
    return "http://" + "localhost" + ":" + port
}
