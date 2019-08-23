package templating

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/redhat-nfvpe/helm2go-operator-sdk/internal/resourcecache"
)

// ReconcileTemplateConfig contains the necessary fields to render the reconcile code block specific to a resource
type ReconcileTemplateConfig struct {
	ResourceImportPackage   string
	ResourceType            string
	ResourceName            string
	LowerResourceName       string
	Kind                    string
	LowerKind               string
	OwnerAPIVersion         string
	PluralLowerResourceName string
}

// NewReconcileTemplateConfig returns a new reconcile templating configuration object
func NewReconcileTemplateConfig(kind, lowerKind string, r *resourcecache.Resource) *ReconcileTemplateConfig {

	resourceType := getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)

	return &ReconcileTemplateConfig{
		Kind:                    kind,
		LowerKind:               lowerKind,
		ResourceImportPackage:   getResourceImportPackage(resourceType),
		ResourceType:            resourceType,
		ResourceName:            getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data),
		LowerResourceName:       strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data)),
		PluralLowerResourceName: getNamePlural(strings.ToLower(getTemplateResourceTitle(&r.GetResourceFunctions()[0].Data))),
	}
}

// Execute renders the template and returns the templated string
func (r *ReconcileTemplateConfig) Execute() (string, error) {

	temp, err := template.New("reconcileBlockTemplate").Parse(r.GetTemplate())
	if err != nil {
		return "", err
	}

	var wr bytes.Buffer
	err = temp.Execute(&wr, r)
	if err != nil {
		return "", err
	}

	return wr.String(), nil
}

// GetTemplate returns the necessary template
func (r *ReconcileTemplateConfig) GetTemplate() string {
	return `
		{{.LowerResourceName}} := {{.PluralLowerResourceName}}.New{{.ResourceName}}ForCR(instance)
		reqLogger.Info({{.LowerResourceName}}.String())
		// Set {{ .Kind }} instance as the owner and controller
		if err{{.ResourceName}} := controllerutil.SetControllerReference(instance, {{.LowerResourceName}}, r.scheme); err{{.ResourceName}} != nil {
			return reconcile.Result{}, err{{.ResourceName}}
		}
		// Check if this {{.ResourceName}} already exists
		found{{.ResourceName}} := &{{.ResourceImportPackage}}.{{.ResourceType}}{}
		err{{.ResourceName}} := r.client.Get(context.TODO(), types.NamespacedName{Name: {{.LowerResourceName}}.Name, Namespace: request.Namespace}, found{{.ResourceName}})
		if err{{.ResourceName}} != nil && errors.IsNotFound(err{{.ResourceName}}) {
			reqLogger.Info("Creating a new {{.ResourceName}}", "{{.LowerResourceName}}.Namespace", request.Namespace, "{{.LowerResourceName}}.Name", {{.LowerResourceName}}.Name)
			{{.LowerResourceName}}.ObjectMeta.SetNamespace(request.Namespace)
			err{{.ResourceName}} = r.client.Create(context.TODO(), {{.LowerResourceName}})
			if err{{.ResourceName}} != nil {
				return reconcile.Result{}, err{{.ResourceName}}
			}
			// Pod created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if err{{.ResourceName}} != nil {
			return reconcile.Result{}, err{{.ResourceName}}
		}
		reqLogger.Info("Skip reconcile: {{.LowerResourceName}} already exists", "{{.ResourceName}}.Namespace",
		found{{.ResourceName}}.Namespace, "svcac{{.LowerResourceName}}cnt.Name", found{{.ResourceName}}.Name)

		if func(instance &{{ .OwnerAPIVersion }}.{{ .Kind }}{}, object &{{.ResourceImportPackage}}.{{.ResourceType}}{}) bool {
			
			// EDIT THIS FILE! tHIS IS SCAFFOLDING FOR YOU TO OWN!
			// The following commented code checks the difference between the 
			// {{ .OwnerAPIVersion }}.{{ .Kind }} instance and the current {{ .ResourceName }} instance.
			// When implementing, you should specify the update conditions.
			
			/*
			var instanceMap map[string]interface{}
			var objectMap map[string]interface{}

			instanceBytes, _ := json.Marshal(instance)
			objectBytes, _ := json.Marshal(object)

			_ = json.Unmarshal(instanceBytes, instanceMap)
			_ = json.Unmarshal(objectBytes, objectMap)

			for objectSpecKey, objectSpec := range objectMap {
				instanceSpec, ok := instanceMap[objectSpecKey]
				if ok {
					if objectSpec != instanceSpec {
						return true
					}
				}
			}
			*/
			return false
		}(instance, found{{.ResourceName}}) {
			
			// The following commented code updates the kubeclient

			// client.Update(found{{.ResourceName}})
			// return reconcile.Result{}, nil
		}
	`
}
