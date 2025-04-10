package content

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestContentMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		want    string
	}{
		{
			name:    "simple text",
			content: FromText("hello"),
			want:    `[{"text":"hello","type":"text"}]`,
		},
		{
			name: "text and image",
			content: FromTextAndImage(
				"hello",
				"https://example.com/image.jpg",
			),
			want: `[{"text":"hello","type":"text"},{"image_url":"https://example.com/image.jpg","type":"imageURL"}]`,
		},
		{
			name:    "json content",
			content: FromRawJSON(json.RawMessage(`{"foo":"bar"}`)),
			want:    `[{"data":{"foo":"bar"},"type":"json"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.content)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestContentUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Content
		wantErr bool
	}{
		{
			name: "simple text",
			json: `[{"type":"text","text":"hello"}]`,
			want: FromText("hello"),
		},
		{
			name: "text and image",
			json: `[
				{"type":"text","text":"hello"},
				{"type":"imageURL","image_url":"https://example.com/image.jpg"}
			]`,
			want: Content{
				&Text{Text: "hello"},
				&ImageURL{URL: "https://example.com/image.jpg"},
			},
		},
		{
			name: "json content",
			json: `[{"type":"json","data":{"foo":"bar"}}]`,
			want: FromRawJSON(json.RawMessage(`{"foo":"bar"}`)),
		},
		{
			name:    "invalid type",
			json:    `[{"type":"invalid"}]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Content
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// For JSON content, we need to compare the marshaled string
			// since the RawMessage might not be byte-for-byte identical
			if len(got) == 1 && got[0].Type() == TypeJSON {
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.want)
				if !reflect.DeepEqual(gotJSON, wantJSON) {
					t.Errorf("UnmarshalJSON() got = %v, want %v", string(gotJSON), string(wantJSON))
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		content Content
	}{
		{
			name:    "simple text",
			content: FromText("hello"),
		},
		{
			name: "text and image",
			content: FromTextAndImage(
				"hello",
				"https://example.com/image.jpg",
			),
		},
		{
			name:    "json content",
			content: FromRawJSON(json.RawMessage(`{"foo":"bar"}`)),
		},
		{
			name: "multiple text items",
			content: Content{
				&Text{Text: "hello"},
				&Text{Text: "world"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.content)
			if err != nil {
				t.Errorf("Marshal error = %v", err)
				return
			}

			var got Content
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Errorf("Unmarshal error = %v", err)
				return
			}

			// For JSON content, compare marshaled results
			if len(tt.content) > 0 && tt.content[0].Type() == TypeJSON {
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.content)
				if !reflect.DeepEqual(gotJSON, wantJSON) {
					t.Errorf("RoundTrip got = %v, want %v", string(gotJSON), string(wantJSON))
				}
				return
			}

			if !reflect.DeepEqual(got, tt.content) {
				t.Errorf("RoundTrip got = %v, want %v", got, tt.content)
			}
		})
	}
}
