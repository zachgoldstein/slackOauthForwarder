# Slack OAUTH Forwarder

## What does this do?
- Executes an oauth request on the Slack API with client id and key from env vars, and a code in the query params 
- Sends an email to an address with all the response data
- Redirects to Polyhack

## To run locally:
```
sam local invoke "slackOauthForwarder" -e testAPIReq.json --debug
```
