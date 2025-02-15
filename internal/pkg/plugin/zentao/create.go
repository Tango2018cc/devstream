package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/goclient"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			goclient.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			goclient.DealWithNsWhenInstall,
			goclient.CreatePersistentVolumeWrapper(defaultPVPath),
			goclient.CreatePersistentVolumeClaim,
			goclient.CreateDeploymentWrapperLabelAndContainerPorts(defaultZentaolabels, &defaultZentaoPorts, defaultName),
			goclient.CreateServiceWrapperLabelAndPorts(defaultZentaolabels, &defaultSVCPort),
			goclient.WaitForReady(retryTimes),
		},
		TerminateOperations: plugininstaller.TerminateOperations{
			goclient.DealWithErrWhenInstall,
		},
		GetStateOperation: goclient.GetState,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(plugininstaller.RawOptions(options))
	if err != nil {
		return nil, err
	}

	// Currently just store zentao's application status: "running" or "stopped"
	return status, nil
}
