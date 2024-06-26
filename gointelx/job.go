package gointelx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/khulnasoft/go-intelx/constants"
)

// UserDetails represents user details in an IntelX job.
type UserDetails struct {
	Username string `json:"username"`
}

// Report represents a report generated by an IntelX job.
type Report struct {
	Name                 string                 `json:"name"`
	Status               string                 `json:"status"`
	Report               map[string]interface{} `json:"report"`
	Errors               []string               `json:"errors"`
	ProcessTime          float64                `json:"process_time"`
	StartTime            time.Time              `json:"start_time"`
	EndTime              time.Time              `json:"end_time"`
	RuntimeConfiguration map[string]interface{} `json:"runtime_configuration"`
	Type                 string                 `json:"type"`
}

// BaseJob respresents all the common fields in a Job and JobList.
type BaseJob struct {
	ID                       int         `json:"id"`
	User                     UserDetails `json:"user"`
	Tags                     []Tag       `json:"tags"`
	ProcessTime              float64     `json:"process_time"`
	IsSample                 bool        `json:"is_sample"`
	Md5                      string      `json:"md5"`
	ObservableName           string      `json:"observable_name"`
	ObservableClassification string      `json:"observable_classification"`
	FileName                 string      `json:"file_name"`
	FileMimetype             string      `json:"file_mimetype"`
	Status                   string      `json:"status"`
	AnalyzersRequested       []string    `json:"analyzers_requested" `
	ConnectorsRequested      []string    `json:"connectors_requested"`
	AnalyzersToExecute       []string    `json:"analyzers_to_execute"`
	ConnectorsToExecute      []string    `json:"connectors_to_execute"`
	ReceivedRequestTime      *time.Time  `json:"received_request_time"`
	FinishedAnalysisTime     *time.Time  `json:"finished_analysis_time"`
	Tlp                      string      `json:"tlp"`
	Errors                   []string    `json:"errors"`
}

// Job represents a job that is being processed in IntelX.
type Job struct {
	BaseJob
	AnalyzerReports  []Report               `json:"analyzer_reports"`
	ConnectorReports []Report               `json:"connector_reports"`
	Permission       map[string]interface{} `json:"permission"`
}

// JobList represents a list of jobs in IntelX.
type JobList struct {
	BaseJob
}

type JobListResponse struct {
	Count      int       `json:"count"`
	TotalPages int       `json:"total_pages"`
	Results    []JobList `json:"results"`
}

// JobService handles communication with job related methods of IntelX API.
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs
type JobService struct {
	client *IntelXClient
}

// List fetches all the jobs in your IntelX instance.
//
//	Endpoint: GET /api/jobs
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_list
func (jobService *JobService) List(ctx context.Context) (*JobListResponse, error) {
	requestUrl := jobService.client.options.Url + constants.BASE_JOB_URL
	contentType := "application/json"
	method := "GET"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return nil, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	jobList := JobListResponse{}
	marashalError := json.Unmarshal(successResp.Data, &jobList)
	if marashalError != nil {
		return nil, marashalError
	}

	return &jobList, nil
}

// Get fetches a specific job through its job ID.
//
//	Endpoint: GET /api/jobs/{jobID}
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_retrieve
func (jobService *JobService) Get(ctx context.Context, jobId uint64) (*Job, error) {
	route := jobService.client.options.Url + constants.SPECIFIC_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId)
	contentType := "application/json"
	method := "GET"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return nil, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	jobResponse := Job{}
	unmarshalError := json.Unmarshal(successResp.Data, &jobResponse)
	if unmarshalError != nil {
		return nil, unmarshalError
	}
	return &jobResponse, nil
}

// DownloadSample fetches the File sample with the given job through its job ID.
//
//	Endpoint: GET /api/jobs/{jobID}/download_sample
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_download_sample_retrieve
func (jobService *JobService) DownloadSample(ctx context.Context, jobId uint64) ([]byte, error) {
	route := jobService.client.options.Url + constants.DOWNLOAD_SAMPLE_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId)
	contentType := "application/json"
	method := "GET"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return nil, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return successResp.Data, nil
}

// Delete removes the given job from your IntelX instance.
//
//	Endpoint: DELETE /api/jobs/{jobID}
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_destroy
func (jobService *JobService) Delete(ctx context.Context, jobId uint64) (bool, error) {
	route := jobService.client.options.Url + constants.SPECIFIC_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId)
	contentType := "application/json"
	method := "DELETE"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

// Kill lets you stop a running job through its ID
//
//	Endpoint: PATCH /api/jobs/{jobID}/kill
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_kill_partial_update
func (jobService *JobService) Kill(ctx context.Context, jobId uint64) (bool, error) {
	route := jobService.client.options.Url + constants.KILL_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId)
	contentType := "application/json"
	method := "PATCH"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

// KillAnalyzer lets you stop an analyzer from running on a processed job through its ID and analyzer name.
//
//	Endpoint: PATCH /api/jobs/{jobID}/analyzer/{nameOfAnalyzer}/kill
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_analyzer_kill_partial_update
func (jobService *JobService) KillAnalyzer(ctx context.Context, jobId uint64, analyzerName string) (bool, error) {
	route := jobService.client.options.Url + constants.KILL_ANALYZER_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId, analyzerName)
	contentType := "application/json"
	method := "PATCH"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

// RetryAnalyzer lets you re-run the selected analyzer on a processed job through its ID and the analyzer name.
//
//	Endpoint: PATCH /api/jobs/{jobID}/analyzer/{nameOfAnalyzer}/retry
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_analyzer_retry_partial_update
func (jobService *JobService) RetryAnalyzer(ctx context.Context, jobId uint64, analyzerName string) (bool, error) {
	route := jobService.client.options.Url + constants.RETRY_ANALYZER_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId, analyzerName)
	contentType := "application/json"
	method := "PATCH"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

// KillConnector lets you stop a connector from running on a processed job through its ID and connector name.
//
//	Endpoint: PATCH /api/jobs/{jobID}/connector/{nameOfConnector}/kill
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_connector_kill_partial_update
func (jobService *JobService) KillConnector(ctx context.Context, jobId uint64, connectorName string) (bool, error) {
	route := jobService.client.options.Url + constants.KILL_CONNECTOR_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId, connectorName)
	contentType := "application/json"
	method := "PATCH"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}

// RetryConnector lets you re-run the selected connector on a processed job through its ID and connector name
//
//	Endpoint: PATCH /api/jobs/{jobID}/connector/{nameOfConnector}/retry
//
// IntelX REST API docs: https://intelx.readthedocs.io/en/latest/Redoc.html#tag/jobs/operation/jobs_connector_retry_partial_update
func (jobService *JobService) RetryConnector(ctx context.Context, jobId uint64, connectorName string) (bool, error) {
	route := jobService.client.options.Url + constants.RETRY_CONNECTOR_JOB_URL
	requestUrl := fmt.Sprintf(route, jobId, connectorName)
	contentType := "application/json"
	method := "PATCH"
	request, err := jobService.client.buildRequest(ctx, method, contentType, nil, requestUrl)
	if err != nil {
		return false, err
	}
	successResp, err := jobService.client.newRequest(ctx, request)
	if err != nil {
		return false, err
	}
	if successResp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	return false, nil
}