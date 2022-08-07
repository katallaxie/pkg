package snowflake

import (
	"encoding/base64"
	"errors"
	"strconv"
	"sync"
	"time"
)

//
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+
//
//

type Node interface {
	Generate() ID
}

const (
	Epoch int64 = 1288834974657
)

var (
	NodeBits uint8 = 10
	StepBits uint8 = 12
)

type node struct {
	epoch time.Time
	time  int64
	id    int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8

	sync.Mutex
}

type ID int64

func (i ID) Int64() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i ID) Bytes() []byte {
	return []byte(i.String())
}

func (i ID) Base64() string {
	return base64.StdEncoding.EncodeToString(i.Bytes())
}

func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}

	return ParseBytes(b)
}

func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

func (i ID) ProtoMessage() *Snowflake {
	return &Snowflake{
		Id: i.Int64(),
	}
}

func (s *Snowflake) ID() ID {
	return ID(s.GetId())
}

func ParseInt64(id int64) ID {
	return ID(id)
}

func New(id int64) (Node, error) {
	n := new(node)
	n.id = id

	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.id < 0 || n.id > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	now := time.Now()
	n.epoch = now.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(now))

	return n, nil
}

func (n *node) Generate() ID {
	n.Lock()
	defer n.Unlock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000

	n.step = 0

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	}

	n.time = now

	r := ID((now)<<n.timeShift |
		(n.id << n.nodeShift) |
		(n.step),
	)

	return r
}
