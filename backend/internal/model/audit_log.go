package model

type AuditLog struct {
    Method         string
    Endpoint       string
    ResourceType   string
    ResourceID     string
    RequestBody    string
    ResponseStatus int
    ResponseBody   string
}