package hex

import (
	"encoding/json"
	"testing"
)

func TestParseHex(t *testing.T) {

	s := `[
		"ac-f7-b0-92-cb-51-11-9f-62-93-d0-8b-2e-06-ab-0b",
		"1b-fc-48-99-21-7d-3b-a9-66-c2-65-80-79-5a-95-3a",
		"ff-5d-97-31-53-c5-87-48-27-37-b8-fc-e8-3a-af-42",
		"a0-22-24-c9-a3-50-ab-c4-e3-35-e9-15-27-11-7e-d6",
		"b3-8d-d7-53-48-79-e0-be-4d-b3-67-e7-1f-46-3d-89",
		"d7-eb-13-4d-31-ee-07-f9-6a-2a-d4-3c-65-07-14-cf",
		"cf-69-86-66-6b-fe-e6-a7-f9-70-a4-16-94-fc-6b-73",
		"33-4b-0c-77-df-a0-56-24-2a-71-08-9d-9e-4b-9a-4d",
		"a9-c3-52-d9-2f-87-f0-03-61-02-1b-78-a2-7c-4e-29",
		"3b-81-0a-13-79-92-65-a4-5b-7a-4f-94-a7-61-41-92",
		"5d-1d-bd-57-20-03-48-d0-54-2b-5a-a6-ac-19-1d-80",
		"c9-5e-8f-de-d5-eb-e9-9c-0e-f4-7c-c7-d2-9d-02-ac",
		"f2-0a-34-64-39-b7-ae-7d-d1-59-fe-6c-f7-7b-b5-aa",
		"7e-a6-9c-cd-25-70-10-47-d4-d0-ac-cc-2d-c4-3f-bb",
		"f1-f5-cc-6a-05-cf-fd-60-96-97-0a-33-d2-e5-59-dd",
		"ce-cb-b4-66-5c-52-eb-33-31-5f-c3-01-a5-12-bb"
	]`

	var q []string

	err := json.Unmarshal([]byte(s), &q)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := json.MarshalIndent(&q, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(b))
}
