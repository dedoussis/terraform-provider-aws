package aws

import (
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
		Update: resourceAwsIotLoggingUpdate,
		Delete: resourceAwsIotLoggingDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"log_level": {
				Type:     schema.TypeString,
				Required: true,
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
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAwsIotLoggingCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn
	params := &iot.SetV2LoggingOptionsInput{
		RoleArn:         aws.String(d.Get("role_arn").(string)),
		DefaultLogLevel: aws.String(d.Get("log_level").(string)),
	}

	if v, ok := d.GetOk("disable_all_longs"); ok {
		params.DisableAllLogs = aws.String(v.(bool))
	}

	log.Printf("[DEBUG] Setting IoT Logging Options: %s", params)
	out, err := conn.SetV2LoggingOptions(params)
	if err != nil {
		return err
	}

	d.SetId(*out.ThingName)

	return resourceAwsIotLoggingRead(d, meta)
}

func resourceAwsIotLoggingRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).iotconn

	params := &iot.GetV2LoggingOptionsInput{}
	log.Printf("[DEBUG] Reading IoT Logging Options: %s", params)
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

func resourceAwsIotLoggingUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsIotLoggingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
