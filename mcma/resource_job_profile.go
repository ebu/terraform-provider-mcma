package mcma

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"time"

	mcmamodel "github.com/ebu/mcma-libraries-go/model"
)

func resourceJobProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Job profile data registered in an MCMA Service Registry",

		CreateContext: resourceJobProfileCreate,
		ReadContext:   resourceJobProfileRead,
		UpdateContext: resourceJobProfileUpdate,
		DeleteContext: resourceJobProfileDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "The MCMA type of resource. This value will always be 'JobProfile'.",
				Computed:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the job profile. MCMA IDs are always absolute urls.",
				Computed:    true,
			},
			"date_created": {
				Type:        schema.TypeString,
				Description: "The date and time at which the job profile data was created.",
				Computed:    true,
			},
			"date_modified": {
				Type:        schema.TypeString,
				Description: "The date and time at which the job profile data was last modified.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the job profile.",
				Required:    true,
			},
			"input_parameter": {
				Type:        schema.TypeSet,
				Description: "A list of input parameters (name and type) that must be provided when running a job for this profile.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the input parameter.",
							Required:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "The type of the input parameter. Should specify an MCMA resource or primitive type.",
							Required:    true,
						},
						"optional": {
							Type:        schema.TypeBool,
							Description: "Flag indicating if this input parameter must be provided or not",
							Optional:    true,
							Default:     false,
						},
					},
				},
			},
			"output_parameter": {
				Type:        schema.TypeSet,
				Description: "A list of output parameters (name and type) that will be set on the job when the service has finished.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the output parameter.",
							Required:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "The type of the output parameter. Should specify an MCMA resource or primitive type.",
							Required:    true,
						},
					},
				},
			},
			"custom_properties": {
				Type:        schema.TypeMap,
				Description: "A collection of key-value pairs specifying additional properties for the job profile.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getJobProfileFromResourceData(d *schema.ResourceData) mcmamodel.JobProfile {
	var inputParameters []mcmamodel.JobParameter
	var optionalInputParameters []mcmamodel.JobParameter
	for _, raw := range d.Get("input_parameter").(*schema.Set).List() {
		pMap := raw.(map[string]interface{})
		inputParameter := mcmamodel.JobParameter{
			ParameterType: pMap["type"].(string),
			ParameterName: pMap["name"].(string),
		}
		if pMap["optional"].(bool) {
			optionalInputParameters = append(optionalInputParameters, inputParameter)
		} else {
			inputParameters = append(inputParameters, inputParameter)
		}
	}

	var outputParameters []mcmamodel.JobParameter
	for _, p := range d.Get("output_parameter").(*schema.Set).List() {
		outputParameter := p.(map[string]interface{})
		outputParameters = append(outputParameters, mcmamodel.JobParameter{
			ParameterType: outputParameter["type"].(string),
			ParameterName: outputParameter["name"].(string),
		})
	}

	return mcmamodel.JobProfile{
		Type:                    "JobProfile",
		Name:                    d.Get("name").(string),
		InputParameters:         inputParameters,
		OutputParameters:        outputParameters,
		OptionalInputParameters: optionalInputParameters,
		CustomProperties:        d.Get("custom_properties").(map[string]interface{}),
	}
}

func resourceJobProfileRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	println("resourceJobProfileRead")
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	jobProfileId := d.Id()
	resource, err := resourceManager.Get(reflect.TypeOf(mcmamodel.JobProfile{}), jobProfileId)
	if err != nil {
		return diag.Errorf("error getting job profile with id %s: %s", jobProfileId, err)
	}
	if resource == nil {
		return diag.Diagnostics{}
	}

	jobProfile := resource.(mcmamodel.JobProfile)

	_ = d.Set("type", jobProfile.Type)
	_ = d.Set("id", jobProfile.Id)
	_ = d.Set("date_created", jobProfile.DateCreated.Format(time.RFC3339))
	_ = d.Set("date_modified", jobProfile.DateModified.Format(time.RFC3339))
	_ = d.Set("name", jobProfile.Name)
	_ = d.Set("custom_properties", jobProfile.CustomProperties)

	var inputParameters []map[string]interface{}
	for _, inputParameter := range jobProfile.InputParameters {
		p := make(map[string]interface{})
		p["type"] = inputParameter.ParameterType
		p["name"] = inputParameter.ParameterName
		p["optional"] = false
		inputParameters = append(inputParameters, p)
	}
	for _, optionalInputParameter := range jobProfile.OptionalInputParameters {
		p := make(map[string]interface{})
		p["type"] = optionalInputParameter.ParameterType
		p["name"] = optionalInputParameter.ParameterName
		p["optional"] = true
		inputParameters = append(inputParameters, p)
	}
	if err = d.Set("input_parameter", inputParameters); err != nil {
		return diag.Errorf("error setting input_parameter for job profile with id %s: %s", jobProfileId, err)
	}

	var outputParameters []map[string]interface{}
	for _, outputParameter := range jobProfile.OutputParameters {
		p := make(map[string]interface{})
		p["type"] = outputParameter.ParameterType
		p["name"] = outputParameter.ParameterName
		outputParameters = append(outputParameters, p)
	}
	if err = d.Set("output_parameter", outputParameters); err != nil {
		return diag.Errorf("error setting output_parameter for job profile with id %s: %s", jobProfileId, err)
	}

	return diag.Diagnostics{}
}

func resourceJobProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	println("resourceJobProfileCreate")
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}
	jobProfile := getJobProfileFromResourceData(d)
	createdResource, err := resourceManager.Create(jobProfile)
	if err != nil {
		return diag.FromErr(err)
	}
	createdJobProfile := createdResource.(mcmamodel.JobProfile)

	d.SetId(createdJobProfile.Id)

	return resourceJobProfileRead(ctx, d, m)
}

func resourceJobProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	println("resourceJobProfileUpdate")
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	jobProfile := getJobProfileFromResourceData(d)
	jobProfile.Id = d.Id()

	_, err := resourceManager.Update(jobProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceJobProfileRead(ctx, d, m)
}

func resourceJobProfileDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	println("resourceJobProfileDelete")
	resourceManager, di := getResourceManager(m)
	if di != nil {
		return di
	}

	err := resourceManager.Delete(reflect.TypeOf(mcmamodel.JobProfile{}), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
