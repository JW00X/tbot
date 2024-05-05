package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"golang.org/x/net/html"
	"strings"	
)

// BongaCamsChecker implements a checker for BongaCams
type BongaCamsChecker struct{ CheckerCommon }

var _ Checker = &BongaCamsChecker{}

type bongacamsModel struct {
	IsAvailable      		bool 	`json:"isAvailable"`
	IsOffline      			bool 	`json:"isOffline"`
	IsPrivatChat      		bool 	`json:"isPrivatChat"`
	IsFullPrivatChat      	bool 	`json:"isFullPrivatChat"`
	IsGroupPrivatChat      	bool 	`json:"isGroupPrivatChat"`
	IsVipShow      			bool 	`json:"isVipShow"`
	DisplayName      		string 	`json:"displayName"`
	RtAvailable      		bool 	`json:"rtAvailable"`
	IsQoQWinner      		bool 	`json:"isQoQWinner"`
	//ProfileImages			string 	`json:"profileImage"`
}

// CheckStatusSingle checks BongaCams model status
func (c *BongaCamsChecker) CheckStatusSingle(modelID string) StatusKind {
	client := c.clientsLoop.nextClient()
	onlineModels = map[string]StatusKind{}
	
	resp, buf, err := onlineQuery(endpoint, client, c.Headers)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot send a query, %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("query status, %d", resp.StatusCode)
	}

	doc, err := html.Parse(bytes.NewReader(buf.Bytes()))
    if err != nil {
        return nil, nil, fmt.Errorf("cannot parse html, %v", err)
    }

	var processAll func(*html.Node, *[]bongacamsModel)
    processAll = func(n *html.Node, ptrModelArr *[]bongacamsModel) {
        if n.Type == html.ElementNode && n.Data == "script" {
            for _, a := range n.Attr {
				if a.Key == "data-type" && strings.Contains(a.Val, "initialState") {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.TextNode && strings.Contains(c.Data, "chatShowStatusOptions") {
							var jsonData map[string]*json.RawMessage
							if err := json.Unmarshal([]byte(c.Data), &jsonData); err != nil {
								Ldbg("cannot parse JSON, %v", err)
							}
							var m bongacamsModel
							err := json.Unmarshal(*jsonData["chatShowStatusOptions"], &m)
							if err != nil {
								Ldbg("cannot parse JSON, %v", err)
							}
							*ptrModelArr = append(*ptrModelArr, m)
						}
					}
				}
			}
 
        }

		for c := n.FirstChild; c != nil; c = c.NextSibling {
            processAll(c, ptrModelArr)
        }
    }

	var models []bongacamsModel
	processAll(doc, &models)

	if len(models) == 0 {
		return nil, nil, errors.New("spec models are not defined")
	}

	for _, m := range models {
		if m.IsOffline == false {
			modelID := strings.ToLower(m.DisplayName)
			onlineModels[modelID] = StatusOnline
			if m.IsPrivatChat == true {
				onlineModels[modelID] = StatusPrivatChat
			}
			if m.IsFullPrivatChat == true {
				onlineModels[modelID] = StatusFullPrivatChat
			}
			if m.IsGroupPrivatChat == true {
				onlineModels[modelID] = StatusGroupPrivatChat
			}
			images[modelID] = "https:" // + m.ProfileImages.ThumbnailImageMediumLive

			return onlineModels[modelID]
		}

		return StatusOffline
	} 
}

// checkEndpoint returns BongaCams online models on the endpoint
func (c *BongaCamsChecker) checkEndpoint(endpoint string) (onlineModels map[string]StatusKind, images map[string]string, err error) {
	client := c.clientsLoop.nextClient()
	onlineModels = map[string]StatusKind{}
	images = map[string]string{}

	log.Printf("[ENDPOINT] %s", endpoint)

	resp, buf, err := onlineQuery(endpoint, client, c.Headers)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot send a query, %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("query status, %d", resp.StatusCode)
	}

	doc, err := html.Parse(bytes.NewReader(buf.Bytes()))
    if err != nil {
        return nil, nil, fmt.Errorf("cannot parse html, %v", err)
    }

	var processAll func(*html.Node, *[]bongacamsModel)
    processAll = func(n *html.Node, ptrModelArr *[]bongacamsModel) {
        if n.Type == html.ElementNode && n.Data == "script" {
            for _, a := range n.Attr {
				if a.Key == "data-type" && strings.Contains(a.Val, "initialState") {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.TextNode && strings.Contains(c.Data, "chatShowStatusOptions") {
							var jsonData map[string]*json.RawMessage
							if err := json.Unmarshal([]byte(c.Data), &jsonData); err != nil {
								Ldbg("cannot parse JSON, %v", err)
							}
							var m bongacamsModel
							err := json.Unmarshal(*jsonData["chatShowStatusOptions"], &m)
							if err != nil {
								Ldbg("cannot parse JSON, %v", err)
							}
							*ptrModelArr = append(*ptrModelArr, m)
						}
					}
				}
			}
 
        }

		for c := n.FirstChild; c != nil; c = c.NextSibling {
            processAll(c, ptrModelArr)
        }
    }

	var models []bongacamsModel
	processAll(doc, &models)

	if len(models) == 0 {
		return nil, nil, errors.New("spec models are not defined")
	}

	for _, m := range models {
		if m.IsOffline == false {
			modelID := strings.ToLower(m.DisplayName)
			onlineModels[modelID] = StatusOnline
			if m.IsPrivatChat == true {
				onlineModels[modelID] = StatusPrivatChat
			}
			if m.IsFullPrivatChat == true {
				onlineModels[modelID] = StatusFullPrivatChat
			}
			if m.IsGroupPrivatChat == true {
				onlineModels[modelID] = StatusGroupPrivatChat
			}
			images[modelID] = "https:" // + m.ProfileImages.ThumbnailImageMediumLive
		}
	}
	return
}

// CheckStatusesMany returns BongaCams online models
func (c *BongaCamsChecker) CheckStatusesMany(QueryModelList, CheckMode) (onlineModels map[string]StatusKind, images map[string]string, err error) {
	return checkEndpoints(c, c.UsersOnlineEndpoints, c.Dbg)
}

// Start starts a daemon
func (c *BongaCamsChecker) Start()                 { c.startFullCheckerDaemon(c) }
func (c *BongaCamsChecker) createUpdater() Updater { return c.createFullUpdater(c) }
