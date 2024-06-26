package lib

import (
	"log"
	"time"
)

// Updater implements updater adapter
type Updater interface {
	QueryUpdates(updateRequest StatusUpdateRequest) error
}

// StatusUpdateRequest represents a request of models updates
type StatusUpdateRequest struct {
	SpecialModels map[string]bool
	Subscriptions map[string]StatusKind
	Specific      map[string]bool
	Callback      func(StatusUpdateResults)
}

// StatusUpdate represents an update of model status
type StatusUpdate struct {
	ModelID string
	Status  StatusKind
}

// StatusUpdateResultsData contains data from updates checking algorithm
type StatusUpdateResultsData struct {
	Updates []StatusUpdate
	Images  map[string]string
	Elapsed time.Duration
}

// StatusUpdateResults contains results from updates checking algorithm
type StatusUpdateResults struct {
	Data   *StatusUpdateResultsData
	Errors int
}

func getUpdates(prev, next map[string]bool, ss map[string]StatusKind) []StatusUpdate {
	var result []StatusUpdate
	log.Printf("[DEBUG] getUpdates: prev: %v next %v", prev, next)
	newElems, removed := HashDiffNewRemoved(prev, next)
	log.Printf("[DEBUG] getUpdates: newElems: %v removed %v", newElems, removed)
	for _, i := range removed {
		if _, ok := ss[i]; ok {
			result = append(result, StatusUpdate{ModelID: i, Status: ss[i]})
		} else {
			result = append(result, StatusUpdate{ModelID: i, Status: StatusOffline})
		}
	}
	for _, i := range newElems {
		result = append(result, StatusUpdate{ModelID: i, Status: StatusOnline})
	}
	log.Printf("[DEBUG] getUpdates: result: %v", result)
	return result
}

func fullUpdateReqToStatus(r StatusUpdateRequest, callback func(StatusResults)) StatusRequest {
	var specific map[string]bool
	if r.Specific != nil {
		specific = map[string]bool{}
		for k := range r.Specific {
			specific[k] = true
		}
	}
	return StatusRequest{
		SpecialModels: r.SpecialModels,
		Specific:      specific,
		Callback:      callback,
	}
}

func selectiveUpdateReqToStatus(r StatusUpdateRequest, callback func(StatusResults)) StatusRequest {
	subs := map[string]bool{}
	for k := range r.Subscriptions {
		subs[k] = true
	}
	for k := range r.Specific {
		subs[k] = true
	}
	return StatusRequest{
		SpecialModels: r.SpecialModels,
		Specific:      subs,
		Callback:      callback,
	}
}

func onlyOnline(ss map[string]StatusKind) map[string]bool {
	boolMap := map[string]bool{}
	for k, s := range ss {
		if s == StatusOnline {
			boolMap[k] = true
		}
	}
	return boolMap
}

func onlineStatuses(ss map[string]bool) map[string]StatusKind {
	boolMap := map[string]StatusKind{}
	for k, s := range ss {
		if s {
			boolMap[k] = StatusOnline
		}
	}
	return boolMap
}
