package task

type ConfigureImplant struct {
	NewJitter         float64
	NewDwell          string
	NewC2CommProtocol string
	ID                uint `json:"id"`
}

// TODO: implement same spec in server
func NewConfigureImplantTask(ID uint, config map[string]interface{}) *ConfigureImplant {

	configureImplant := new(ConfigureImplant)
	if raw, ok := config["NewDwell"]; ok {
		if s, ok := raw.(string); ok {
			configureImplant.NewDwell = s
		}
	}
	if raw, ok := config["NewJitter"]; ok {
		if f, ok := raw.(float64); ok {
			configureImplant.NewJitter = f
		}
	}
	// TODO: honestly just using string everywhere just get rid of ENUM?
	if raw, ok := config["NewC2CommProtocol"]; ok {
		if s, ok := raw.(string); ok {
			configureImplant.NewC2CommProtocol = s
		}
	}
	configureImplant.ID = ID
	return configureImplant
}

// THIS WILL NOT BE CALLED AS CONFIGURE IMPLANT is a special task
func (p *ConfigureImplant) Do() Result {
	return *new(Result)
}

func (p *ConfigureImplant) GetID() uint {
	return p.ID
}
