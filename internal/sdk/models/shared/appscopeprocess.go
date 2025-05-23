// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type Interface struct {
	Config    *string `json:"config,omitempty"`
	Connected *bool   `json:"connected,omitempty"`
	Name      *string `json:"name,omitempty"`
}

func (o *Interface) GetConfig() *string {
	if o == nil {
		return nil
	}
	return o.Config
}

func (o *Interface) GetConnected() *bool {
	if o == nil {
		return nil
	}
	return o.Connected
}

func (o *Interface) GetName() *string {
	if o == nil {
		return nil
	}
	return o.Name
}

type AppScopeProcessProcess struct {
	HostPid   *float64 `json:"hostPid,omitempty"`
	ID        *string  `json:"id,omitempty"`
	MachineID *string  `json:"machine_id,omitempty"`
	Pid       float64  `json:"pid"`
	UUID      *string  `json:"uuid,omitempty"`
}

func (o *AppScopeProcessProcess) GetHostPid() *float64 {
	if o == nil {
		return nil
	}
	return o.HostPid
}

func (o *AppScopeProcessProcess) GetID() *string {
	if o == nil {
		return nil
	}
	return o.ID
}

func (o *AppScopeProcessProcess) GetMachineID() *string {
	if o == nil {
		return nil
	}
	return o.MachineID
}

func (o *AppScopeProcessProcess) GetPid() float64 {
	if o == nil {
		return 0.0
	}
	return o.Pid
}

func (o *AppScopeProcessProcess) GetUUID() *string {
	if o == nil {
		return nil
	}
	return o.UUID
}

type AppScopeProcess struct {
	Cfg              *AppscopeConfigWithCustom `json:"cfg,omitempty"`
	ConfigID         *string                   `json:"config_id,omitempty"`
	ID               string                    `json:"id"`
	Interfaces       []Interface               `json:"interfaces,omitempty"`
	LastError        *string                   `json:"lastError,omitempty"`
	Process          *AppScopeProcessProcess   `json:"process,omitempty"`
	ProcessingStatus *AppScopeProcessingStatus `json:"processingStatus,omitempty"`
	SourceID         *string                   `json:"source_id,omitempty"`
	Status           *AppScopeProcessStatus    `json:"status,omitempty"`
}

func (o *AppScopeProcess) GetCfg() *AppscopeConfigWithCustom {
	if o == nil {
		return nil
	}
	return o.Cfg
}

func (o *AppScopeProcess) GetConfigID() *string {
	if o == nil {
		return nil
	}
	return o.ConfigID
}

func (o *AppScopeProcess) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *AppScopeProcess) GetInterfaces() []Interface {
	if o == nil {
		return nil
	}
	return o.Interfaces
}

func (o *AppScopeProcess) GetLastError() *string {
	if o == nil {
		return nil
	}
	return o.LastError
}

func (o *AppScopeProcess) GetProcess() *AppScopeProcessProcess {
	if o == nil {
		return nil
	}
	return o.Process
}

func (o *AppScopeProcess) GetProcessingStatus() *AppScopeProcessingStatus {
	if o == nil {
		return nil
	}
	return o.ProcessingStatus
}

func (o *AppScopeProcess) GetSourceID() *string {
	if o == nil {
		return nil
	}
	return o.SourceID
}

func (o *AppScopeProcess) GetStatus() *AppScopeProcessStatus {
	if o == nil {
		return nil
	}
	return o.Status
}
