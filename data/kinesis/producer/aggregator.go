package producer

import (
	"bytes"
	"crypto/md5"

	k "github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/monetha/mth-core/data/kinesis/producer/messages"
	"google.golang.org/protobuf/proto"
)

var (
	magicNumber = []byte{0xF3, 0x89, 0x9A, 0xC2}
)

// Aggregator used to aggregate records into kinesis.PutRecordsRequestEntry
type Aggregator struct {
	buf    []*messages.Record
	pkeys  []string
	nbytes int
}

// Size return how many bytes stored in the aggregator.
// including partition keys.
func (a *Aggregator) Size() int {
	return a.nbytes
}

// Count return how many records stored in the aggregator.
func (a *Aggregator) Count() int {
	return len(a.buf)
}

// Put record using `data` and `partitionKey`. This method is not thread-safe.
func (a *Aggregator) Put(data []byte, partitionKey string) {
	// For now, all records in the aggregated record will have
	// the same partition key.
	// later, we will add shard-mapper same as the KPL use.
	// see: https://github.com/a8m/kinesis-producer/issues/1
	if len(a.pkeys) == 0 {
		a.pkeys = append(a.pkeys, partitionKey)
		a.nbytes += len([]byte(partitionKey))
	}
	keyIndex := uint64(len(a.pkeys) - 1)
	a.buf = append(a.buf, &messages.Record{
		Data:              data,
		PartitionKeyIndex: &keyIndex,
	})
	a.nbytes += len(data)
}

// Drain create an aggregated `kinesis.PutRecordsRequestEntry`
// that compatible with the KCL's deaggregation logic.
//
// If you interested to know more about it. see: aggregation-format.md
func (a *Aggregator) Drain() (*k.PutRecordsRequestEntry, error) {
	data, err := proto.Marshal(&messages.AggregatedRecord{
		PartitionKeyTable: a.pkeys,
		Records:           a.buf,
	})
	if err != nil {
		return nil, err
	}
	h := md5.New()
	h.Write(data)
	checkSum := h.Sum(nil)
	aggData := append(magicNumber, data...)
	aggData = append(aggData, checkSum...)
	entry := &k.PutRecordsRequestEntry{
		Data:         aggData,
		PartitionKey: &a.pkeys[0],
	}
	a.clear()
	return entry, nil
}

func (a *Aggregator) clear() {
	a.buf = make([]*messages.Record, 0)
	a.pkeys = make([]string, 0)
	a.nbytes = 0
}

// Test if a given entry is aggregated record.
func isAggregated(entry *k.PutRecordsRequestEntry) bool {
	return bytes.HasPrefix(entry.Data, magicNumber)
}

func extractRecords(entry *k.PutRecordsRequestEntry) (out []*k.PutRecordsRequestEntry) {
	src := entry.Data[len(magicNumber) : len(entry.Data)-md5.Size]
	dest := new(messages.AggregatedRecord)
	err := proto.Unmarshal(src, dest)
	if err != nil {
		return
	}
	for i := range dest.Records {
		r := dest.Records[i]
		out = append(out, &k.PutRecordsRequestEntry{
			Data:         r.GetData(),
			PartitionKey: &dest.PartitionKeyTable[r.GetPartitionKeyIndex()],
		})
	}
	return
}
