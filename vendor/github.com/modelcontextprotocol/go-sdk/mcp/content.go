// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package mcp

import (
	"encoding/json"
	"errors"
	"fmt"
)

// A Content is a [TextContent], [ImageContent], [AudioContent] or
// [EmbeddedResource].
//
// TODO(rfindley): add ResourceLink.
type Content interface {
	MarshalJSON() ([]byte, error)
	fromWire(*wireContent)
}

// TextContent is a textual content.
type TextContent struct {
	Text        string
	Meta        Meta
	Annotations *Annotations
}

func (c *TextContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireContent{
		Type:        "text",
		Text:        c.Text,
		Meta:        c.Meta,
		Annotations: c.Annotations,
	})
}

func (c *TextContent) fromWire(wire *wireContent) {
	c.Text = wire.Text
	c.Meta = wire.Meta
	c.Annotations = wire.Annotations
}

// ImageContent contains base64-encoded image data.
type ImageContent struct {
	Meta        Meta
	Annotations *Annotations
	Data        []byte // base64-encoded
	MIMEType    string
}

func (c *ImageContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireContent{
		Type:        "image",
		MIMEType:    c.MIMEType,
		Data:        c.Data,
		Meta:        c.Meta,
		Annotations: c.Annotations,
	})
}

func (c *ImageContent) fromWire(wire *wireContent) {
	c.MIMEType = wire.MIMEType
	c.Data = wire.Data
	c.Meta = wire.Meta
	c.Annotations = wire.Annotations
}

// AudioContent contains base64-encoded audio data.
type AudioContent struct {
	Data        []byte
	MIMEType    string
	Meta        Meta
	Annotations *Annotations
}

func (c AudioContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireContent{
		Type:        "audio",
		MIMEType:    c.MIMEType,
		Data:        c.Data,
		Meta:        c.Meta,
		Annotations: c.Annotations,
	})
}

func (c *AudioContent) fromWire(wire *wireContent) {
	c.MIMEType = wire.MIMEType
	c.Data = wire.Data
	c.Meta = wire.Meta
	c.Annotations = wire.Annotations
}

// EmbeddedResource contains embedded resources.
type EmbeddedResource struct {
	Resource    *ResourceContents
	Meta        Meta
	Annotations *Annotations
}

func (c *EmbeddedResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(&wireContent{
		Type:        "resource",
		Resource:    c.Resource,
		Meta:        c.Meta,
		Annotations: c.Annotations,
	})
}

func (c *EmbeddedResource) fromWire(wire *wireContent) {
	c.Resource = wire.Resource
	c.Meta = wire.Meta
	c.Annotations = wire.Annotations
}

// ResourceContents contains the contents of a specific resource or
// sub-resource.
type ResourceContents struct {
	URI      string `json:"uri"`
	MIMEType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     []byte `json:"blob,omitempty"`
	Meta     Meta   `json:"_meta,omitempty"`
}

func (r ResourceContents) MarshalJSON() ([]byte, error) {
	// If we could assume Go 1.24, we could use omitzero for Blob and avoid this method.
	if r.URI == "" {
		return nil, errors.New("ResourceContents missing URI")
	}
	if r.Blob == nil {
		// Text. Marshal normally.
		type wireResourceContents ResourceContents // (lacks MarshalJSON method)
		return json.Marshal((wireResourceContents)(r))
	}
	// Blob.
	if r.Text != "" {
		return nil, errors.New("ResourceContents has non-zero Text and Blob fields")
	}
	// r.Blob may be the empty slice, so marshal with an alternative definition.
	br := struct {
		URI      string `json:"uri,omitempty"`
		MIMEType string `json:"mimeType,omitempty"`
		Blob     []byte `json:"blob"`
	}{
		URI:      r.URI,
		MIMEType: r.MIMEType,
		Blob:     r.Blob,
	}
	return json.Marshal(br)
}

// wireContent is the wire format for content.
// It represents the protocol types TextContent, ImageContent, AudioContent
// and EmbeddedResource.
// The Type field distinguishes them. In the protocol, each type has a constant
// value for the field.
// At most one of Text, Data, and Resource is non-zero.
type wireContent struct {
	Type        string            `json:"type"`
	Text        string            `json:"text,omitempty"`
	MIMEType    string            `json:"mimeType,omitempty"`
	Data        []byte            `json:"data,omitempty"`
	Resource    *ResourceContents `json:"resource,omitempty"`
	Meta        Meta              `json:"_meta,omitempty"`
	Annotations *Annotations      `json:"annotations,omitempty"`
}

func contentsFromWire(wires []*wireContent, allow map[string]bool) ([]Content, error) {
	var blocks []Content
	for _, wire := range wires {
		block, err := contentFromWire(wire, allow)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func contentFromWire(wire *wireContent, allow map[string]bool) (Content, error) {
	if allow != nil && !allow[wire.Type] {
		return nil, fmt.Errorf("invalid content type %q", wire.Type)
	}
	switch wire.Type {
	case "text":
		v := new(TextContent)
		v.fromWire(wire)
		return v, nil
	case "image":
		v := new(ImageContent)
		v.fromWire(wire)
		return v, nil
	case "audio":
		v := new(AudioContent)
		v.fromWire(wire)
		return v, nil
	case "resource":
		v := new(EmbeddedResource)
		v.fromWire(wire)
		return v, nil
	}
	return nil, fmt.Errorf("internal error: unrecognized content type %s", wire.Type)
}
