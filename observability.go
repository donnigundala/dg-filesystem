package dgfilesystem

import (
	"context"
	"io"
	"time"

	"github.com/donnigundala/dg-core/contracts/filesystem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrumentationName = "github.com/donnigundala/dg-filesystem"
)

// ObservedDisk wraps a filesystem.Disk and records metrics.
type ObservedDisk struct {
	disk              filesystem.Disk
	name              string
	requestCounter    metric.Int64Counter
	durationHistogram metric.Float64Histogram
	bytesCounter      metric.Int64Counter
}

// NewObservedDisk creates a new ObservedDisk decorator.
func NewObservedDisk(disk filesystem.Disk, name string) *ObservedDisk {
	meter := otel.GetMeterProvider().Meter(instrumentationName)

	requestCounter, _ := meter.Int64Counter(
		"filesystem.disk.operation.count",
		metric.WithDescription("Total number of disk operations"),
	)

	durationHistogram, _ := meter.Float64Histogram(
		"filesystem.disk.operation.duration",
		metric.WithDescription("Duration of disk operations"),
		metric.WithUnit("ms"),
	)

	bytesCounter, _ := meter.Int64Counter(
		"filesystem.disk.operation.bytes",
		metric.WithDescription("Total bytes read or written"),
		metric.WithUnit("By"),
	)

	return &ObservedDisk{
		disk:              disk,
		name:              name,
		requestCounter:    requestCounter,
		durationHistogram: durationHistogram,
		bytesCounter:      bytesCounter,
	}
}

func (d *ObservedDisk) record(operation string, startTime time.Time, err error, bytes int64, byteType string) {
	duration := float64(time.Since(startTime).Milliseconds())
	status := "success"
	if err != nil {
		status = "error"
	}

	attrs := metric.WithAttributes(
		attribute.String("disk.name", d.name),
		attribute.String("operation", operation),
		attribute.String("status", status),
	)

	d.requestCounter.Add(context.Background(), 1, attrs)
	d.durationHistogram.Record(context.Background(), duration, attrs)

	if bytes > 0 {
		byteAttrs := metric.WithAttributes(
			attribute.String("disk.name", d.name),
			attribute.String("operation", operation),
			attribute.String("type", byteType),
		)
		d.bytesCounter.Add(context.Background(), bytes, byteAttrs)
	}
}

func (d *ObservedDisk) Put(path string, content []byte) error {
	start := time.Now()
	err := d.disk.Put(path, content)
	d.record("put", start, err, int64(len(content)), "write")
	return err
}

func (d *ObservedDisk) PutStream(path string, content io.Reader) error {
	start := time.Now()
	err := d.disk.PutStream(path, content)
	// We don't easily know the byte count for streams without wrapping the reader
	d.record("put_stream", start, err, 0, "write")
	return err
}

func (d *ObservedDisk) Get(path string) ([]byte, error) {
	start := time.Now()
	content, err := d.disk.Get(path)
	d.record("get", start, err, int64(len(content)), "read")
	return content, err
}

func (d *ObservedDisk) GetStream(path string) (io.ReadCloser, error) {
	start := time.Now()
	reader, err := d.disk.GetStream(path)
	d.record("get_stream", start, err, 0, "read")
	return reader, err
}

func (d *ObservedDisk) Exists(path string) (bool, error) {
	start := time.Now()
	exists, err := d.disk.Exists(path)
	d.record("exists", start, err, 0, "")
	return exists, err
}

func (d *ObservedDisk) Delete(path string) error {
	start := time.Now()
	err := d.disk.Delete(path)
	d.record("delete", start, err, 0, "")
	return err
}

func (d *ObservedDisk) Url(path string) string {
	return d.disk.Url(path)
}

func (d *ObservedDisk) SignedUrl(path string, expiration time.Duration) (string, error) {
	start := time.Now()
	url, err := d.disk.SignedUrl(path, expiration)
	d.record("signed_url", start, err, 0, "")
	return url, err
}

func (d *ObservedDisk) MakeDirectory(path string) error {
	start := time.Now()
	err := d.disk.MakeDirectory(path)
	d.record("make_directory", start, err, 0, "")
	return err
}

func (d *ObservedDisk) DeleteDirectory(path string) error {
	start := time.Now()
	err := d.disk.DeleteDirectory(path)
	d.record("delete_directory", start, err, 0, "")
	return err
}
