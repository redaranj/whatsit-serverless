# whatsit-serverless

This project provides a set of AWS Lambda functions for sending WhatsApp messages.

**Do not use for sensitive/private information.**

Supported use cases:

* dispatching [alertmanager](https://github.com/prometheus/alertmanager) alerts
* dispatching [datadog](https://datadoghq.com) alerts

## Setup

#### Prerequisites

* go 1.12
* node 8.10
* [serverless](https://github.com/serverless/serverless)

```

## Configure

The service expects to find this parameter in the AWS SecretsManager as `SecureString`s

* `/whatsit-serverless/<STAGE>/api-key`

Where `<STAGE>` is the `development`/`production` stage you specify.

You can use the AWS cli to add them, or terraform, as you like.

As part of the serverless deployment, the functions will be given an IAM policy
that allows them read access to these params.

## Deploy

#### Deploy development

```
make deploy-dev
```

#### Deploy production

```
make deploy-prod
```

### Remove infrastructure

```
serverless remove --stage XXX
```

## Usage Example with alertmanager

### 1. Sign in to WhatsApp

Signing in to WhatsApp requires you to scan a QR code from the WhatsApp Settings screen on your phone. The process is as follows:

Send a request to the `register` endpoint with the number you want to send messages from and your email address:

```bash
curl -s  -XPOST -d '{"number": "15555555555", email: "test@example.com" }' https://<LAMBDA_ENDPOINT>/register?api_key=<API_KEY>
```

You will receive an email with a QR code. Scan it within 20 seconds to complete the sign in process.

At that point the registration command will return a numberId (SHA256 of your phone number) and a secret. You will use these in subsequent requests.

### 2. Verify the number is registered

```bash
curl -s  -XPOST -d '{"number": "15555555555" }' https://<LAMBDA_ENDPOINT>/verify?secret=<SECRET>
```

### 3. Send test alert

```bash
curl -s  -XPOST -d '{"sender": <YOUR_NUMBER_ID>, "number": "15555555555", "message": "Hi" }' https://<LAMBDA_ENDPOINT>/alert?secret=<SECRET>
```

### 4. Configure alertmanager

In your `alertmanager.yml` add this block to your `receivers:` section

```yaml
receivers:
  - name: 'my-serverless-whatsapp'
    webhook_configs:
      - url: 'https://<LAMBDA_ENDPOINT>/alert?secret=<SECRET>'
```

### Delete a number

To delete a number, and its stored state (session keys, etc):

```bash
curl -s  -XPOST -d '{"sender": <YOUR_NUMBER_ID>}' https://<LAMBDA_ENDPOINT>/delete?secret=<SECRET>
```

## HTTP API Reference

### Authentication

Authentication is handled via either an API key (for registration) or a per-WhatsApp-number shared secret.

The `register` endpoint requires a query parameter `api_key` that must equal the value configured for the `WHATSIT_API_KEY` variable.

* If the API key is missing or blank on the server, the server will return a `401`.
* If the API key is missing from the client's request, the server will return a `401`.

All other endpoints require a query parameter `secret` which will be returned by the registration endpoint.

* If this secret is missing or blank on the server, the server will return a `401`.
* If this secret is missing from the client's request, the server will return a `401`.

### Number IDs

In all lifecycle endpoints (register, verify, delete), the Signal account
number (used to send messages) is passed as the fully qualified number (e.g.,
`15555555555`).

However in the endpoints that send messages, i.e., `/alert`, the raw number  is
not used. Instead a `sender` query parameter is required.  The value of
`sender` is the sha256 hash of the fully qualified number. This hash is
returned as part of the response from `/verify` and `/register`, so if you note
that down, there is no need to calculate it yourself.

This is done in order to prevent third-party services (alertmanager, datadog,
etc) from knowing what your Signal number is.

### POST /register

**JSON Payload:**

* `number`: string containing the fully qualified number, with no spaces or punctuation
* `email`: an email address where you can receive a QR code

**Example Request:**
```json
{
    "number": "15555555555",
    "email": "test@example.com"
}
```

**Response Code:**

* `200`: if the registration request succeeded
* `400`: if the request payload is invalid
* `500`: if the  registration request failed

**Response Body:**

* `result`: a result message
* `numberId`: the id of the number registered (the sha256 hash), to be used when sending messages
* `secret`: the secret to include in subsequent requests

**Example Response:**

```json
{
    "result": "registration complete",
    "numberId": "910a625c4ba147b544e6bd2f267e130ae14c591b6ba9c25cb8573322dedbebd0",
    "secret": "xxx"
}
```

### POST /verify

**JSON Payload:**

* `number`: string containing the fully qualified number, with no spaces or punctuation

**Example Request:**

```json
{
    "number": "15555555555"
}
```

**Response Code:**

* `200`: if the verification succeeded
* `400`: if the request payload is invalid
* `500`: if the verification failed

**Response Body:**

* `result`: a result message
* `numberId`: the id of the number registered (the sha256 hash), to be used when sending messages

**Example Response:**

```json
{
    "result": "the number '15555555555' was previously registered and can send messages",
    "numberId": "910a625c4ba147b544e6bd2f267e130ae14c591b6ba9c25cb8573322dedbebd0"
}
```

### POST /delete

To delete a number from the backend (but not from signal's servers).

**JSON Payload:**

* `number`: string containing the fully qualified number, with no spaces or punctuation

**Example Request:**

```json
{
    "number": "15555555555"
}
```

**Response Code:**

* `200`: if the deletion succeeded
* `400`: if the request payload is invalid
* `500`: if the deletion failed

**Response Body:**

* `result`: a result message

**Example Response:**

```json
{
    "result":"the number '15555555555' is deleted"
}
```

### POST /alert

An endpoint that receives [alertmanager]() formatted json payloads containing alerts. The Signal recipients are decided based on the
`receiver` attribute of the payload. See

**JSON Payload:**

See [alertmanager's <webhook_config>](https://prometheus.io/docs/alerting/configuration/#webhook_config)


**Example Request:**



**Response Code:**

* `200`: if the send succeeded
* `400`: if the request payload is invalid, or if the `sender` query parameter is missing
* `404`: if the `sender` query parameter is invalid (refers to non-existent number in the backend)
* `500`: if the send failed

**Response Body:**

* `result`: a result message

**Example Response:**

```json
{
    "result": "ok"
}
```

# License

`whatsit-serverless` is licensed under the [GNU Affero General Public License
(AGPL) v3+](https://www.gnu.org/licenses/agpl-3.0.en.html).

Copyright (C) 2019 Darren Clarke <darren@redaranj.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
