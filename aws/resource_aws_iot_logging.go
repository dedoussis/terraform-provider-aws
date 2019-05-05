package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAwsIotLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsIotLoggingCreate,
		Read:   resourceAwsIotLoggingRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"log_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEBUG",
					"INFO",
					"ERROR",
					"WARN",
					"DISABLED",
				}, false),
			},
			"disable_all_logs": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAwsIotLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	roleArn := d.Get("role_arn").(string)
	logLevel := d.Get("log_level").(string)
	disableAllLogs := d.Get("disable_all_longs").(bool)

	params := &iot.SetV2LoggingOptionsInput{
		RoleArn:         aws.String(roleArn),
		DefaultLogLevel: aws.String(logLevel),
		DisableAllLogs:  aws.Bool(disableAllLogs),
	}

	log.Printf("[DEBUG] Setting IoT Logging Options: %s", params)
	_, err := conn.SetV2LoggingOptions(params)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s|%s", roleArn, logLevel))
	return resourceAwsIotLoggingRead(d, meta)
}

func resourceAwsIotLoggingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	params := &iot.GetV2LoggingOptionsInput{}
	log.Printf("[DEBUG] Retrieving IoT Logging Options: %s", params)
	out, err := conn.GetV2LoggingOptions(params)

	if err != nil {
		if isAWSErr(err, iot.ErrCodeResourceNotFoundException, "") {
			log.Printf("[WARN] IoT Logging Options %q not found, removing from state", d.Id())
			d.SetId("")
		}
		return err
	}

	log.Printf("[DEBUG] Retrieved IoT Logging Options: %s", out)

	d.Set("role_arn", out.RoleArn)
	d.Set("default_log_level", out.DefaultLogLevel)
	d.Set("disable_all_logs", out.DisableAllLogs)
	return nil
}
