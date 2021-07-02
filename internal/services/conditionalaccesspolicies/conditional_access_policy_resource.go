package conditionalaccesspolicies

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func conditionalAccessPolicyResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: conditionalAccessPolicyResourceCreate,
		ReadContext:   conditionalAccessPolicyResourceRead,
		UpdateContext: conditionalAccessPolicyResourceUpdate,
		DeleteContext: conditionalAccessPolicyResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"state": {
				Type:     schema.TypeString,
				Required: true,
			},

			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"applications": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"included_applications": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_applications": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"included_user_actions": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"users": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"included_users": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_users": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"included_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"included_roles": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_roles": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"client_app_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"locations": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"included_locations": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_locations": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"platforms": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"included_platforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"excluded_platforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"sign_in_risk_levels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_risk_levels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"grant_controls": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},

						"built_in_controls": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"custom_authentication_factors": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"terms_of_use": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"session_controls": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_enforced_restrictions": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"cloud_app_security": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"cloud_app_security_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"persistent_browser": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"sign_in_frequency": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func conditionalAccessPolicyResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ConditionalAccessPolicies.MsClient

	displayName := d.Get("display_name").(string)
	state := d.Get("state").(string)

	conditionsRaw := d.Get("conditions").([]interface{})
	conditions := expandConditionalAccessConditionSet(conditionsRaw)

	grantControlsRaw := d.Get("grant_controls").([]interface{})
	grantControls := expandConditionalAccessGrantControls(grantControlsRaw)

	sessionControlsRaw := d.Get("session_controls").([]interface{})
	sessionControls := expandConditionalAccessSessionControls(sessionControlsRaw)

	properties := msgraph.ConditionalAccessPolicy{
		DisplayName:     utils.String(displayName),
		State:           utils.String(state),
		Conditions:      conditions,
		GrantControls:   grantControls,
		SessionControls: sessionControls,
	}

	policy, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Could not create conditional access policy")
	}

	if policy.ID == nil || *policy.ID == "" {
		return tf.ErrorDiagF(errors.New("Bad API response"), "Object ID returned for conditional access policy is nil/empty")
	}

	d.SetId(*policy.ID)

	return conditionalAccessPolicyResourceRead(ctx, d, meta)
}

func conditionalAccessPolicyResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ConditionalAccessPolicies.MsClient

	properties := msgraph.ConditionalAccessPolicy{
		ID: utils.String(d.Id()),
	}

	if d.HasChange("display_name") {
		displayName := d.Get("display_name").(string)
		properties.DisplayName = &displayName
	}

	if d.HasChange("state") {
		state := d.Get("state").(string)
		properties.State = &state
	}

	if d.HasChange("conditions") {
		conditionsRaw := d.Get("conditions").([]interface{})
		conditions := expandConditionalAccessConditionSet(conditionsRaw)
		properties.Conditions = conditions
	}

	if d.HasChange("grant_controls") {
		grantControlsRaw := d.Get("grant_controls").([]interface{})
		grantControls := expandConditionalAccessGrantControls(grantControlsRaw)
		properties.GrantControls = grantControls
	}

	if d.HasChange("session_controls") {
		sessionControlsRaw := d.Get("session_controls").([]interface{})
		sessionControls := expandConditionalAccessSessionControls(sessionControlsRaw)
		properties.SessionControls = sessionControls
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Could not update conditional access policy with ID: %q", d.Id())
	}

	return nil
}

func conditionalAccessPolicyResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ConditionalAccessPolicies.MsClient

	policy, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Conditional Access Policy with Object ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving Conditional Access Policy with object ID %q", d.Id())
	}

	tf.Set(d, "display_name", policy.DisplayName)
	tf.Set(d, "id", policy.ID)
	tf.Set(d, "state", policy.State)
	tf.Set(d, "conditions", flattenConditionalAccessConditionSet(policy.Conditions))
	tf.Set(d, "grant_controls", flattenConditionalAccessGrantControls(policy.GrantControls))
	tf.Set(d, "session_controls", flattenConditionalAccessSessionControls(policy.SessionControls))

	return nil
}

func conditionalAccessPolicyResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ConditionalAccessPolicies.MsClient

	_, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Conditional Access Policy with ID %q already deleted", d.Id())
			return nil
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving conditional access policy with ID %q", d.Id())
	}

	status, err = client.Delete(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting conditional access policy with ID %q, got status %d", d.Id(), status)
	}

	return nil
}

func flattenConditionalAccessConditionSet(in *msgraph.ConditionalAccessConditionSet) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"applications":        flattenConditionalAccessApplications(in.Applications),
			"users":               flattenConditionalAccessUsers(in.Users),
			"client_app_types":    tf.FlattenStringSlicePtr(in.ClientAppTypes),
			"locations":           flattenConditionalAccessLocations(in.Locations),
			"platforms":           flattenConditionalAccessPlatforms(in.Platforms),
			"sign_in_risk_levels": tf.FlattenStringSlicePtr(in.SignInRiskLevels),
			"user_risk_levels":    tf.FlattenStringSlicePtr(in.UserRiskLevels),
		},
	}
}

func flattenConditionalAccessApplications(in *msgraph.ConditionalAccessApplications) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"included_applications": tf.FlattenStringSlicePtr(in.IncludeApplications),
			"excluded_applications": tf.FlattenStringSlicePtr(in.ExcludeApplications),
			"included_user_actions": tf.FlattenStringSlicePtr(in.IncludeUserActions),
		},
	}
}

func flattenConditionalAccessUsers(in *msgraph.ConditionalAccessUsers) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"included_users":  tf.FlattenStringSlicePtr(in.IncludeUsers),
			"excluded_users":  tf.FlattenStringSlicePtr(in.ExcludeUsers),
			"included_groups": tf.FlattenStringSlicePtr(in.IncludeGroups),
			"excluded_groups": tf.FlattenStringSlicePtr(in.ExcludeGroups),
			"included_roles":  tf.FlattenStringSlicePtr(in.IncludeRoles),
			"excluded_roles":  tf.FlattenStringSlicePtr(in.ExcludeRoles),
		},
	}
}

func flattenConditionalAccessLocations(in *msgraph.ConditionalAccessLocations) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"included_locations": tf.FlattenStringSlicePtr(in.IncludeLocations),
			"excluded_locations": tf.FlattenStringSlicePtr(in.ExcludeLocations),
		},
	}
}

func flattenConditionalAccessPlatforms(in *msgraph.ConditionalAccessPlatforms) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"included_platforms": tf.FlattenStringSlicePtr(in.IncludePlatforms),
			"excluded_platforms": tf.FlattenStringSlicePtr(in.ExcludePlatforms),
		},
	}
}

func flattenConditionalAccessGrantControls(in *msgraph.ConditionalAccessGrantControls) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"operator":                      in.Operator,
			"built_in_controls":             tf.FlattenStringSlicePtr(in.BuiltInControls),
			"custom_authentication_factors": tf.FlattenStringSlicePtr(in.CustomAuthenticationFactors),
			"terms_of_use":                  tf.FlattenStringSlicePtr(in.TermsOfUse),
		},
	}
}

func flattenConditionalAccessSessionControls(in *msgraph.ConditionalAccessSessionControls) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"application_enforced_restrictions": flattenApplicationEnforcedRestrictionsSessionControl(in.ApplicationEnforcedRestrictions),
			"cloud_app_security":                flattenCloudAppSecurityControl(in.CloudAppSecurity),
			"persistent_browser":                flattenPersistentBrowserSessionControl(in.PersistentBrowser),
			"sign_in_frequency":                 flattenSignInFrequencySessionControl(in.SignInFrequency),
		},
	}
}

func flattenApplicationEnforcedRestrictionsSessionControl(in *msgraph.ApplicationEnforcedRestrictionsSessionControl) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": in.IsEnabled,
		},
	}
}

func flattenCloudAppSecurityControl(in *msgraph.CloudAppSecurityControl) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                 in.IsEnabled,
			"cloud_app_security_type": in.CloudAppSecurityType,
		},
	}
}

func flattenPersistentBrowserSessionControl(in *msgraph.PersistentBrowserSessionControl) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": in.IsEnabled,
			"mode":    in.Mode,
		},
	}
}

func flattenSignInFrequencySessionControl(in *msgraph.SignInFrequencySessionControl) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": in.IsEnabled,
			"type":    in.Type,
			"value":   in.Value,
		},
	}
}

func expandConditionalAccessConditionSet(in []interface{}) *msgraph.ConditionalAccessConditionSet {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessConditionSet{}
	config := in[0].(map[string]interface{})

	applications := config["applications"].([]interface{})
	users := config["users"].([]interface{})
	clientAppTypes := config["client_app_types"].([]interface{})
	locations := config["locations"].([]interface{})
	platforms := config["platforms"].([]interface{})
	signInRiskLevels := config["sign_in_risk_levels"].([]interface{})
	userRiskLevels := config["user_risk_levels"].([]interface{})

	result.Applications = expandConditionalAccessApplications(applications)
	result.Users = expandConditionalAccessUsers(users)
	result.ClientAppTypes = tf.ExpandStringSlicePtr(clientAppTypes)
	result.Locations = expandConditionalAccessLocations(locations)
	result.Platforms = expandConditionalAccessPlatforms(platforms)
	result.SignInRiskLevels = tf.ExpandStringSlicePtr(signInRiskLevels)
	result.UserRiskLevels = tf.ExpandStringSlicePtr(userRiskLevels)

	return &result
}

func expandConditionalAccessApplications(in []interface{}) *msgraph.ConditionalAccessApplications {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessApplications{}
	config := in[0].(map[string]interface{})

	includeApplications := config["included_applications"].([]interface{})
	excludeApplications := config["excluded_applications"].([]interface{})
	includeUserActions := config["included_user_actions"].([]interface{})

	result.IncludeApplications = tf.ExpandStringSlicePtr(includeApplications)
	result.ExcludeApplications = tf.ExpandStringSlicePtr(excludeApplications)
	result.IncludeUserActions = tf.ExpandStringSlicePtr(includeUserActions)

	return &result
}

func expandConditionalAccessUsers(in []interface{}) *msgraph.ConditionalAccessUsers {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessUsers{}
	config := in[0].(map[string]interface{})

	includeUsers := config["included_users"].([]interface{})
	excludeUsers := config["excluded_users"].([]interface{})
	includeGroups := config["included_groups"].([]interface{})
	excludeGroups := config["excluded_groups"].([]interface{})
	includeRoles := config["included_roles"].([]interface{})
	excludeRoles := config["excluded_roles"].([]interface{})

	result.IncludeUsers = tf.ExpandStringSlicePtr(includeUsers)
	result.ExcludeUsers = tf.ExpandStringSlicePtr(excludeUsers)
	result.IncludeGroups = tf.ExpandStringSlicePtr(includeGroups)
	result.ExcludeGroups = tf.ExpandStringSlicePtr(excludeGroups)
	result.IncludeRoles = tf.ExpandStringSlicePtr(includeRoles)
	result.ExcludeRoles = tf.ExpandStringSlicePtr(excludeRoles)

	return &result
}

func expandConditionalAccessPlatforms(in []interface{}) *msgraph.ConditionalAccessPlatforms {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessPlatforms{}
	config := in[0].(map[string]interface{})

	includePlatforms := config["included_platforms"].([]interface{})
	excludePlatforms := config["excluded_platforms"].([]interface{})

	result.IncludePlatforms = tf.ExpandStringSlicePtr(includePlatforms)
	result.ExcludePlatforms = tf.ExpandStringSlicePtr(excludePlatforms)

	return &result
}

func expandConditionalAccessLocations(in []interface{}) *msgraph.ConditionalAccessLocations {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessLocations{}
	config := in[0].(map[string]interface{})

	includeLocations := config["included_locations"].([]interface{})
	excludeLocations := config["excluded_locations"].([]interface{})

	result.IncludeLocations = tf.ExpandStringSlicePtr(includeLocations)
	result.ExcludeLocations = tf.ExpandStringSlicePtr(excludeLocations)

	return &result
}

func expandConditionalAccessGrantControls(in []interface{}) *msgraph.ConditionalAccessGrantControls {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessGrantControls{}
	config := in[0].(map[string]interface{})

	operator := config["operator"].(string)
	builtInControls := config["built_in_controls"].([]interface{})
	customAuthenticationFactors := config["custom_authentication_factors"].([]interface{})
	termsOfUse := config["terms_of_use"].([]interface{})

	result.Operator = &operator
	result.BuiltInControls = tf.ExpandStringSlicePtr(builtInControls)
	result.CustomAuthenticationFactors = tf.ExpandStringSlicePtr(customAuthenticationFactors)
	result.TermsOfUse = tf.ExpandStringSlicePtr(termsOfUse)

	return &result
}

func expandConditionalAccessSessionControls(in []interface{}) *msgraph.ConditionalAccessSessionControls {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ConditionalAccessSessionControls{}
	config := in[0].(map[string]interface{})

	applicationEnforcedRestrictions := config["application_enforced_restrictions"].([]interface{})
	cloudAppSecurity := config["cloud_app_security"].([]interface{})
	persistentBrowser := config["persistent_browser"].([]interface{})
	signInFrequency := config["sign_in_frequency"].([]interface{})

	result.ApplicationEnforcedRestrictions = expandApplicationEnforcedRestrictionsSessionControl(applicationEnforcedRestrictions)
	result.CloudAppSecurity = expandCloudAppSecurityControl(cloudAppSecurity)
	result.PersistentBrowser = expandPersistentBrowserSessionControl(persistentBrowser)
	result.SignInFrequency = expandSignInFrequencySessionControl(signInFrequency)

	return &result
}

func expandApplicationEnforcedRestrictionsSessionControl(in []interface{}) *msgraph.ApplicationEnforcedRestrictionsSessionControl {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.ApplicationEnforcedRestrictionsSessionControl{}
	config := in[0].(map[string]interface{})

	enabled := config["enabled"].(bool)

	result.IsEnabled = &enabled

	return &result
}

func expandCloudAppSecurityControl(in []interface{}) *msgraph.CloudAppSecurityControl {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.CloudAppSecurityControl{}
	config := in[0].(map[string]interface{})

	enabled := config["enabled"].(bool)
	cloudAppSecurityType := config["cloud_app_security_type"].(string)

	result.IsEnabled = &enabled
	result.CloudAppSecurityType = &cloudAppSecurityType

	return &result
}

func expandPersistentBrowserSessionControl(in []interface{}) *msgraph.PersistentBrowserSessionControl {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.PersistentBrowserSessionControl{}
	config := in[0].(map[string]interface{})

	enabled := config["enabled"].(bool)
	mode := config["mode"].(string)

	result.IsEnabled = &enabled
	result.Mode = &mode

	return &result
}

func expandSignInFrequencySessionControl(in []interface{}) *msgraph.SignInFrequencySessionControl {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.SignInFrequencySessionControl{}
	config := in[0].(map[string]interface{})

	enabled := config["enabled"].(bool)
	controlType := config["type"].(string)
	value := int32(config["value"].(int))

	result.IsEnabled = &enabled
	result.Type = &controlType
	result.Value = &value

	return &result
}
