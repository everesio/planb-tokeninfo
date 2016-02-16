package revoke

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Revocation struct {
	Type      string // token, claim, global
	Data      map[string]interface{}
	Timestamp int
}

type jsonRevoke struct {
	Meta struct {
		ForceRefresh string `json:"force_refresh"`
	} `json:"meta"`
	Revocation []struct {
		Type      string `json:"type"` // TOKEN, CLAIM, GLOBAL
		RevokedAt string `json:"revoked_at"`
		Data      struct {
			Name         string `json:"name,omitempty"`          // CLAIM
			ValueHash    string `json:"value_hash,omitempty"`    // CLAIM
			IssuedBefore string `json:"issued_before,omitempty"` // CLAIM, GLOBAL
			TokenHash    string `json:"token_hash,omitempty"`    // TOKEN
			RevokedAt    string `json:"revoked_at,omitempty"`    // TOKEN
		} `json:"data"`
	} `json:"revocations"`
}

func (r *jsonRevoke) UnmarshallJSON(data []byte, forcedRefresh bool) (err error) {
	var buf jsonRevoke
	if err = json.Unmarshall(data, &buf); err != nil {
		log.Errorf("Error unmarshalling revocation json. " + err.Error())
		return err
	}

	// Note: if we already foreced a refresh, we don't want to do it again
	// otherwise we'll get stuck in an infinite loop
	if buf.ForceRefresh != "" && !forcedRefresh {
		i, err := strconv.Atoi(buf.ForceRefresh)
		if err != nil {
			log.Errorf("Error converting ForceRefresh to int." + err.Error())
		} else {
			// TODO: not sure how to get the current Cache from here. . .
			refreshCacheFromTime(i)
		}
	}

	return
}

func (r []*Revocation) getRevocationFromJson(json *jsonRevoke.Revocation) {

	t := int32(time.Now().Unix())
	for j := range json {

		switch j.Type {
		case "TOKEN":
			valid := isHashTimestampValid(j.Data.TokenHash, j.Data.RevokedAt)
			if !valid {
				log.Errorf("Invalid revocation data. TokenHash: %s, RevokedAt: %s", j.Data.TokenHash, j.Data.RevokedAt)
				continue
			}
			r.Data["token_hash"] = j.Data.TokenHash
			r.Data["revoked_at"] = j.Data.RevokedAt
		case "CLAIM":
			valid := isHashTimestampValid(j.Data.ValueHash, j.Data.IssuedBefore)
			if !valid {
				log.Errorf("Invalid revocation data. ValueHash: %s, IssuedBefore: %s", j.Data.ValueHash, j.Data.IssuedBefore)
				continue
			}
			r.Claims[j.Data.Name][j.Data.ValueHash] = i
			r.Data["value_hash"] = j.Data.ValueHash
			r.Data["issued_before"] = j.Data.IssuedBefore
			r.Data["name"] = j.Data.Name
		case "GLOBAL":
			_, err := strconv.Atoi(j.Data.IssuedBefore)
			if err != nil {
				log.Errorf("Erorr converting IssuedBefore to int. " + err.Error())
				continue
			}
			r.Data["issued_before"] = j.Data.IssuedBefore
		default:
			log.Errorf("Unsupported revocation type: %s", j.Type)
			continue
		}

		r.Type = j.Type
		r.Timestamp = t
	}
	return
}

func isHashTimestampValid(hash, timestamp string) bool {
	if hash == "" || timestamp == "" {
		return false
	}

	_i, err := strconv.Atoi(timestamp)
	if err != nil {
		log.Errorf("Erorr converting timestamp to int. " + err.Error())
		return false
	}

	return true

}
