// Code generated by network-console-controller codegen. DO NOT EDIT
package api

// Implements ResponseSetter and CollectionResponseSetter for the generated
// response objects

// SetCount
func (r *AddressListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *AddressListResponse) SetResults(v []AddressRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *AddressListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *AddressResponse) SetResults(v AddressRecord) {
	r.Results = v
}

// SetCount
func (r *ApplicationFlowResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ApplicationFlowResponse) SetResults(v []ApplicationFlowRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ApplicationFlowResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetCount
func (r *ConnectionListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ConnectionListResponse) SetResults(v []ConnectionRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ConnectionListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetCount
func (r *ConnectorListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ConnectorListResponse) SetResults(v []ConnectorRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ConnectorListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *ConnectorResponse) SetResults(v ConnectorRecord) {
	r.Results = v
}

// SetCount
func (r *FlowAggregateListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *FlowAggregateListResponse) SetResults(v []FlowAggregateRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *FlowAggregateListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *FlowAggregateResponse) SetResults(v FlowAggregateRecord) {
	r.Results = v
}

// SetCount
func (r *LinkListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *LinkListResponse) SetResults(v []LinkRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *LinkListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *LinkResponse) SetResults(v LinkRecord) {
	r.Results = v
}

// SetCount
func (r *ListenerListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ListenerListResponse) SetResults(v []ListenerRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ListenerListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *ListenerResponse) SetResults(v ListenerRecord) {
	r.Results = v
}

// SetCount
func (r *ProcessGroupListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ProcessGroupListResponse) SetResults(v []ProcessGroupRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ProcessGroupListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *ProcessGroupResponse) SetResults(v ProcessGroupRecord) {
	r.Results = v
}

// SetCount
func (r *ProcessListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *ProcessListResponse) SetResults(v []ProcessRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *ProcessListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *ProcessResponse) SetResults(v ProcessRecord) {
	r.Results = v
}

// SetCount
func (r *RouterAccessListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *RouterAccessListResponse) SetResults(v []RouterAccessRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *RouterAccessListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *RouterAccessResponse) SetResults(v RouterAccessRecord) {
	r.Results = v
}

// SetCount
func (r *RouterLinkListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *RouterLinkListResponse) SetResults(v []RouterLinkRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *RouterLinkListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *RouterLinkResponse) SetResults(v RouterLinkRecord) {
	r.Results = v
}

// SetCount
func (r *RouterListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *RouterListResponse) SetResults(v []RouterRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *RouterListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *RouterResponse) SetResults(v RouterRecord) {
	r.Results = v
}

// SetCount
func (r *SiteListResponse) SetCount(v int64) {
	r.Count = v
}

// SetResults
func (r *SiteListResponse) SetResults(v []SiteRecord) {
	r.Results = v
}

// SetTimeRangeCount
func (r *SiteListResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// SetResults
func (r *SiteResponse) SetResults(v SiteRecord) {
	r.Results = v
}

// SetCount
func (r *CollectionResponse) SetCount(v int64) {
	r.Count = v
}

// SetTimeRangeCount
func (r *CollectionResponse) SetTimeRangeCount(v int64) {
	r.TimeRangeCount = v
}

// Implements Record interface for the generated record objects

// GetEndTime
func (r AddressRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r AddressRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ApplicationFlowRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ApplicationFlowRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ConnectionRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ConnectionRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ConnectorRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ConnectorRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r FlowAggregateRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r FlowAggregateRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r LinkRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r LinkRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ListenerRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ListenerRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ProcessGroupRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ProcessGroupRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r ProcessRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r ProcessRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r RouterAccessRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r RouterAccessRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r RouterLinkRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r RouterLinkRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r RouterRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r RouterRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r SiteRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r SiteRecord) GetStartTime() uint64 {
	return r.StartTime
}

// GetEndTime
func (r BaseRecord) GetEndTime() uint64 {
	return r.EndTime
}

// GetStartTime
func (r BaseRecord) GetStartTime() uint64 {
	return r.StartTime
}
