package udp

import (
	"datagram-toolkit/netem"
	"datagram-toolkit/util/mocks"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	logger := stderrLogger
	netemCfg := netem.Config{
		WriteFragmentSize: 48,
		WriteReorderNth:   2,
		WriteDuplicateNth: 3,
		WriteLossNth:      4,
	}
	streamCfg := DefaultStreamConfig()
	streamCfg.Logger = logger

	rand := rand.New(rand.NewSource(0))
	expectedLen := 64
	expected := make([]byte, expectedLen)

	c1, c2 := mocks.Conn()
	s1 := NewStream(netem.New(c1, netemCfg), streamCfg)
	s2 := NewStream(netem.New(c2, netemCfg), streamCfg)

	logger.Printf("s1: %p", s1)
	logger.Printf("s2: %p", s2)

	require := require.New(t)
	_, err := io.ReadFull(rand, expected)
	require.Nil(err)

	w, err := s1.Write(expected)
	require.Nil(err)
	require.Equal(expectedLen, w)

	buf := make([]byte, expectedLen*2)
	r, err := s2.Read(buf)
	require.Nil(err)
	require.Equal(expectedLen, r)
	require.Equal(expected, buf[:r])

	require.Nil(s1.Close())
	require.Equal(io.EOF, s2.Close())
}
