package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// metadata holds all the information about previous generation.
type metadata struct {
	// Input represents the previous input information used to generate the output.
	Input struct {
		Content string // Content is the previous input file content.
		Prompt  string // Prompt is the prompt used in generating the previous output.
	}
	// Output holds additional information about previously generated output.
	Output struct {
		Sha256 string // Sha256 is the hex-encoded SHA-256 checksum of the generated output.
	}
}

type metadataView struct {
	Prompt string
	Data   string // Compressed base64 content
}

type metadataViewData struct {
	Input struct {
		Content string
	}
	Output struct {
		Sha256 string
	}
}

func metadataRead(metaPath string) (*metadata, error) {
	content, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var v metadataView
	if err := yaml.Unmarshal(content, &v); err != nil {
		return nil, err
	}

	zb, err := base64.StdEncoding.DecodeString(string(v.Data))
	if err != nil {
		return nil, err
	}

	gr, err := gzip.NewReader(bytes.NewReader(zb))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gr); err != nil {
		gr.Close()
		return nil, err
	}
	gr.Close()

	var vd metadataViewData

	gd := gob.NewDecoder(&buf)
	if err := gd.Decode(&vd); err != nil {
		return nil, err
	}

	var md metadata

	md.Input.Content = vd.Input.Content
	md.Input.Prompt = v.Prompt
	md.Output.Sha256 = vd.Output.Sha256

	return &md, nil
}

func metadataWrite(metaPath string, md *metadata) error {
	var vd metadataViewData

	vd.Input.Content = md.Input.Content
	vd.Output.Sha256 = md.Output.Sha256

	var gobBuf bytes.Buffer
	ge := gob.NewEncoder(&gobBuf)
	if err := ge.Encode(&vd); err != nil {
		return err
	}

	var gzipBuf bytes.Buffer
	gw := gzip.NewWriter(&gzipBuf)
	if _, err := gw.Write(gobBuf.Bytes()); err != nil {
		gw.Close()
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}

	zb := gzipBuf.Bytes()
	zbuf := make([]byte, base64.StdEncoding.EncodedLen(len(zb)))
	base64.StdEncoding.Encode(zbuf, zb)

	v := metadataView{
		Prompt: md.Input.Prompt,
		Data:   string(zbuf),
	}

	vb, err := yaml.Marshal(&v)
	if err != nil {
		return err
	}

	return atomicWrite(metaPath, vb)
}
