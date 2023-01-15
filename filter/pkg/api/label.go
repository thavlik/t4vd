package api

import (
	"encoding/json"
	"time"
)

func (m *Label) UnmarshalJSON(body []byte) error {
	v := make(map[string]interface{})
	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}
	*m = *NewLabelFromMap(v)
	return nil
}

func (m *Label) MarshalJSON() ([]byte, error) {
	p := make(map[string]interface{})
	p["id"] = m.ID
	p["projectID"] = m.ProjectID
	p["gadgetID"] = m.GadgetID
	p["creatorID"] = m.CreatorID
	p["created"] = m.Created
	if len(m.Payload) > 0 {
		payload := make(map[string]interface{})
		for k, v := range m.Payload {
			payload[k] = v
		}
		p["payload"] = payload
	}
	if len(m.Tags) > 0 {
		p["tags"] = m.Tags
	}
	if m.Parent != nil {
		p["parent"] = m.Parent
	}
	if m.Seek != 0 {
		p["seek"] = int64(m.Seek)
	}
	if m.Pad != 0 {
		p["pad"] = int64(m.Pad)
	}
	return json.Marshal(p)
}

func NewLabelFromMap(p map[string]interface{}) *Label {
	m := &Label{}
	var ok bool
	if m.ID, ok = p["id"].(string); !ok {
		// try mongo _id
		m.ID, _ = p["_id"].(string)
	}
	m.CreatorID, _ = p["creatorID"].(string)
	m.Comment, _ = p["comment"].(string)
	m.GadgetID, _ = p["gadgetID"].(string)
	if v, ok := p["created"].(int64); ok {
		t := time.Unix(0, v)
		m.Created = &t
	}
	m.DeleterID, _ = p["deleterID"].(string)
	if v, ok := p["deleted"].(int64); ok {
		t := time.Unix(0, v)
		m.Deleted = &t
	}
	if parent, ok := p["parent"].(map[string]interface{}); ok {
		m.Parent = NewLabelFromMap(parent)
	}
	if v, ok := p["seek"].(int64); ok {
		m.Seek = time.Duration(v)
	}
	if v, ok := p["pad"].(int64); ok {
		m.Pad = time.Duration(v)
	}
	if v, ok := p["tags"].([]string); ok {
		m.Tags = v
	} else if v, ok := p["tags"].([]interface{}); ok {
		m.Tags = make([]string, len(v))
		for i, vv := range v {
			m.Tags[i], _ = vv.(string)
		}
	}
	m.Payload, _ = p["payload"].(map[string]interface{})
	return m
}
