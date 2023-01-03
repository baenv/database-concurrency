package transducer

type Config struct {
	State    State
	Data     interface{}
	Metadata Metadata
}
type Metadata struct {
	ChildConfig     map[string]*Config
	ParallelConfigs []*Config
	Test            func() error
}

func CreateConfig() *Config {
	return &Config{}
}

func (c *Config) GetState() State {
	return c.State
}

func (c *Config) GetChildState(name string) State {
	if c.Metadata.ChildConfig != nil {
		_, exists := c.Metadata.ChildConfig[name]
		if exists {
			return c.Metadata.ChildConfig[name].State
		}
	}
	return State(Invalid)
}

func (c *Config) SetState(state State) *Config {
	c.State = state
	return c
}

func (c *Config) SetData(data interface{}) *Config {
	c.Data = data
	return c
}

func (c *Config) SetMetadata(metadata Metadata) *Config {
	c.Metadata = metadata
	return c
}

func (c *Config) SetChildConfig(name string, config *Config) *Config {
	if c.Metadata.ChildConfig == nil {
		c.Metadata.ChildConfig = map[string]*Config{}
	}
	_, exists := c.Metadata.ChildConfig[name]
	if !exists {
		c.Metadata.ChildConfig[name] = CreateConfig()
	}
	c.Metadata.ChildConfig[name] = config
	return c
}
