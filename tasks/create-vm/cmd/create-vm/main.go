package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/constants"
	errors2 "github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/errors"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/output"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/utils/parse"
	"github.com/suomiy/kubevirt-tekton-tasks/tasks/create-vm/pkg/vmcreator"
)

var (
	templateName              string
	templateNamespace         string
	templateParams            = parse.NewSpaceSeparatedListFlag(TemplateParamsOptionName)
	dataVolumes               = parse.NewSpaceSeparatedListFlag(DataVolumesOptionName)
	ownDataVolumes            = parse.NewSpaceSeparatedListFlag(OwnDataVolumesOptionName)
	persistentVolumeClaims    = parse.NewSpaceSeparatedListFlag(PersistentVolumeClaimsOptionName)
	ownPersistentVolumeClaims = parse.NewSpaceSeparatedListFlag(OwnPersistentVolumeClaimsOptionName)
	showOutput                string
)

func init() {
	flag.CommandLine.SortFlags = false
	flag.StringVar(&templateName, TemplateNameOptionName, "", "Name of a template to create VM from")
	flag.StringVar(&templateNamespace, TemplateNamespaceOptionName, "", "Namespace of a template to create VM from")
	flag.Var(templateParams, TemplateParamsOptionName, "Template params to pass when processing the template manifest.\nEach param name should be followed by its value.\nEg NAME my-vm DESC blue")
	flag.Var(dataVolumes, DataVolumesOptionName, "Add DVs to VM Volumes.\nEg dv1 dv2")
	flag.Var(ownDataVolumes, OwnDataVolumesOptionName, "Add DVs to VM Volumes and add VM to DV ownerReferences.\nThese DataVolumes will be deleted once the created VM gets deleted.\nEg dv1 dv2")
	flag.Var(persistentVolumeClaims, PersistentVolumeClaimsOptionName, "Add PVCs to VM Volumes.\nEg pvc1 pvc2")
	flag.Var(ownPersistentVolumeClaims, OwnPersistentVolumeClaimsOptionName, "Add PVCs to VM Volumes and add VM to PVC ownerReferences.\nThese PVCs will be deleted once the created VM gets deleted.\nEg pvc1 pvc2")
	flag.StringVar(&showOutput, OutputParamOptionName, "", fmt.Sprintf("Output format. One of: %v", strings.Join(output.GetOutputTypeNames(), "|")))
}

func checkArgErors() error {
	var requiredStrings = map[string]string{
		TemplateNameOptionName:      templateName,
		TemplateNamespaceOptionName: templateNamespace,
	}

	return parse.RequireStringArgs(requiredStrings)
}

// TODO order of flags matters
// SpaceSeparatedListFlags must be at the end!
// deprecate pflag and refactor to use custom parsing all together
func customParse(flags ...*parse.SpaceSeparatedListFlag) {
	for _, f := range flags {
		err := f.SetReal()
		if err != nil {
			utils.Exit(WrongArgsExitCode, err)
		}
	}
}

func main() {
	fmt.Println(os.Args) // TODO remove
	flag.Parse()
	customParse(templateParams, dataVolumes, ownDataVolumes, persistentVolumeClaims, ownPersistentVolumeClaims)

	if err := checkArgErors(); err != nil {
		utils.Exit(WrongArgsExitCode, err)
	}

	cliParams := &parse.CLIParams{
		TemplateName:              templateName,
		TemplateNamespace:         templateNamespace,
		TemplateParams:            templateParams.GetMapValues(),
		DataVolumes:               dataVolumes.GetValues(),
		OwnDataVolumes:            ownDataVolumes.GetValues(),
		PersistentVolumeClaims:    persistentVolumeClaims.GetValues(),
		OwnPersistentVolumeClaims: ownPersistentVolumeClaims.GetValues(),
	}

	vmCreator := vmcreator.NewVMCreator(cliParams)

	if err := vmCreator.CheckVolumesExist(); err != nil {
		utils.Exit(VolumesNotPresentExitCode, err)
	}

	vm, err := vmCreator.CreateVM()

	if err != nil {
		utils.ExitOrDie(CreateVMErrorExitCode, err,
			errors2.IsStatusErrorSoft(err, http.StatusNotFound, http.StatusConflict, http.StatusUnprocessableEntity),
		)
	}

	// write results
	resultsDir := GetTektonResultsDir()
	utils.WriteToFile(filepath.Join(resultsDir, NameResultName), vm.Name)
	utils.WriteToFile(filepath.Join(resultsDir, NamespaceResultName), vm.Namespace)

	if err := vmCreator.OwnVolumes(vm); err != nil {
		utils.Exit(OwnVolumesErrorExitCode, err)
	}

	output.PrettyPrint(vm, output.OutputType(showOutput))
}
