package transducer

type Outputs struct {
	Config  *Config
	Effects []Effect
}

func CreateOutputs() *Outputs {
	return &Outputs{Config: &Config{}}
}

func (o *Outputs) GetState() State {
	if o.Config != nil {
		return o.Config.State
	}
	return State(Invalid)
}

func (o *Outputs) GetChildState(name string) State {
	if o.Config.Metadata.ChildConfig != nil {
		_, exists := o.Config.Metadata.ChildConfig[name]
		if exists {
			return o.Config.Metadata.ChildConfig[name].State
		}
	}
	return State(Invalid)
}

func (o *Outputs) SetState(state State) *Outputs {
	o.Config.SetState(state)
	return o
}

func (o *Outputs) SetData(data interface{}) *Outputs {
	o.Config.Data = data
	return o
}

func (o *Outputs) SetMetadata(metadata Metadata) *Outputs {
	o.Config.Metadata = metadata
	return o
}

func (o *Outputs) AddEffect(effect Effect) *Outputs {
	o.Effects = append(o.Effects, effect)
	return o
}

func (o *Outputs) AddEffects(effects []Effect) *Outputs {
	o.Effects = append(o.Effects, effects...)
	return o
}

func (o *Outputs) TransduceChild(transducer Transducer, config *Config, input Input) *Outputs {
	if config.Metadata.ChildConfig != nil {
		_, exists := config.Metadata.ChildConfig[transducer.Name]
		if exists {
			innerOutputs := transducer.Transduce(config.Metadata.ChildConfig[transducer.Name], input)
			outputs := MergeChildOutputs(o, innerOutputs, transducer.Name)
			return outputs
		}
	}
	return o
}
