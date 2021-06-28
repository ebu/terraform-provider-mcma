package mcma

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/diag"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	mcmaclient "github.com/ebu/mcma-libraries-go/client"
	mcmamodel "github.com/ebu/mcma-libraries-go/model"
)

func resourceJobProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobProfileCreate,
		ReadContext:   resourceJobProfileRead,
		UpdateContext: resourceJobProfileUpdate,
		DeleteContext: resourceJobProfileDelete,

		Schema: map[string]*schema.Schema{
			"@type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"input_parameters": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: false,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Optional: false,
						},
					},
				},
			},
			"output_parameters": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: false,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Optional: false,
						},
					},
				},
			},
			"optional_input_parameters": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_name": {
							Type:     schema.TypeString,
							Optional: false,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Optional: false,
						},
					},
				},
			},
			"custom_properties": {
				Type:     schema.TypeMap,
				Required: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getJobProfileFromResourceData(d *schema.ResourceData) mcmamodel.JobProfile {
	var inputParameters []mcmamodel.JobParameter
	for _, p := range d.Get("input_parameters").([]interface{}) {
		inputParameter := p.(map[string]interface{})
		inputParameters = append(inputParameters, mcmamodel.JobParameter{
			ParameterType: inputParameter["parameter_type"].(string),
			ParameterName: inputParameter["parameter_name"].(string),
		})
	}

	var outputParameters []mcmamodel.JobParameter
	for _, p := range d.Get("output_parameters").([]interface{}) {
		outputParameter := p.(map[string]interface{})
		outputParameters = append(outputParameters, mcmamodel.JobParameter{
			ParameterType: outputParameter["parameter_type"].(string),
			ParameterName: outputParameter["parameter_name"].(string),
		})
	}

	var optionalInputParameters []mcmamodel.JobParameter
	for _, p := range d.Get("optional_input_parameters").([]interface{}) {
		optionalInputParameter := p.(map[string]interface{})
		optionalInputParameters = append(optionalInputParameters, mcmamodel.JobParameter{
			ParameterType: optionalInputParameter["parameter_type"].(string),
			ParameterName: optionalInputParameter["parameter_name"].(string),
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
	resourceManager := m.(*mcmaclient.ResourceManager)

	jobProfileId := d.Id()
	resource, err := resourceManager.Get("JobProfile", jobProfileId)
	if err != nil {
		return diag.Errorf("error getting job profile with id %s: %s", jobProfileId, err)
	}

	jobProfile := resource.(mcmamodel.JobProfile)

	_ = d.Set("type", jobProfile.Type)
	_ = d.Set("id", jobProfile.Id)
	_ = d.Set("date_created", jobProfile.DateCreated)
	_ = d.Set("date_modified", jobProfile.DateModified)
	_ = d.Set("name", jobProfile.Name)
	_ = d.Set("custom_properties", jobProfile.CustomProperties)

	var inputParameters []map[string]interface{}
	for _, inputParameter := range jobProfile.InputParameters {
		p := make(map[string]interface{})
		p["parameter_type"] = inputParameter.ParameterType
		p["parameter_name"] = inputParameter.ParameterName
		inputParameters = append(inputParameters, p)
	}
	if err = d.Set("input_parameters", inputParameters); err != nil {
		return diag.Errorf("error setting input_parameters for job profile with id %s: %s", jobProfileId, err)
	}

	var outputParameters []map[string]interface{}
	for _, outputParameter := range jobProfile.OutputParameters {
		p := make(map[string]interface{})
		p["parameter_type"] = outputParameter.ParameterType
		p["parameter_name"] = outputParameter.ParameterName
		outputParameters = append(outputParameters, p)
	}
	if err = d.Set("output_parameters", outputParameters); err != nil {
		return diag.Errorf("error setting output_parameters for job profile with id %s: %s", jobProfileId, err)
	}

	var optionalInputParameters []map[string]interface{}
	for _, optionalInputParameter := range jobProfile.OptionalInputParameters {
		p := make(map[string]interface{})
		p["parameter_type"] = optionalInputParameter.ParameterType
		p["parameter_name"] = optionalInputParameter.ParameterName
		optionalInputParameters = append(optionalInputParameters, p)
	}
	if err = d.Set("optional_input_parameters", optionalInputParameters); err != nil {
		return diag.Errorf("error setting optional_input_parameters for job profile with id %s: %s", jobProfileId, err)
	}

	return diag.Diagnostics{}
}

func resourceJobProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

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
	resourceManager := m.(*mcmaclient.ResourceManager)

	jobProfile := getJobProfileFromResourceData(d)
	jobProfile.Id = d.Id()

	_, err := resourceManager.Update(jobProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceJobProfileRead(ctx, d, m)
}

func resourceJobProfileDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceManager := m.(*mcmaclient.ResourceManager)

	err := resourceManager.Delete("JobProfile", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
