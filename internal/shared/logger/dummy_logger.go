package logger

type DummyLogger struct{}

func (DummyLogger) Info(string, ...Field)  {}
func (DummyLogger) Error(string, ...Field) {}
func (DummyLogger) Fatal(string, ...Field) {}
func (DummyLogger) Debug(string, ...Field) {}
