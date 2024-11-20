# Notifyre Go Client

## APIs:
- [List Faxes - Click for Demo](./internal/examples/list_faxes.go)
- [Send Fax - Click for Demo](./internal/examples/send_fax.go)

## Webhooks

**A HIGHLY IMPORTANT NOTE ABOUT WEBHOOKS**
*Due to Notifyre NOT returning statuses of failed faxes, a rough approach to capturing this error is used. At a large scale, this implementation will kill. Use wisely, and request that Notifyre add the error codes to the webhook events.*
[For details about why this is a rough implementation, look at the comment here](./fax/attempt_to_get_failed_fax_error.go)

When a FaxSent webhook receives a failure event, it will attempt to lookup the Fax ID in recents. **Be sure to handle "fax ID not found" error**

- [FaxSent Webhook](./internal/examples/fax_sent_webhook.go)
