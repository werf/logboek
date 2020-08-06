package stream

import "github.com/werf/logboek/internal/stream/fitter"

type proxyStream struct {
	*Stream
}

func (s proxyStream) Write(data []byte) (int, error) {
	if s.IsMuted() {
		return 0, nil
	}

	if !s.IsProxyStreamDataFormattingEnabled() {
		return s.logFBase("%s", string(data))
	}

	for _, chunk := range splitData(data, 256) {
		msg := string(chunk)

		if s.Stream.IsLineWrappingEnabled() {
			msg = fitter.FitText(msg, &s.State.State, s.ContentWidth(), true, true)
		}

		_, err := s.processAndLogFBase("%s", msg)
		if err != nil {
			return len(data), err
		}
	}

	return len(data), nil
}
