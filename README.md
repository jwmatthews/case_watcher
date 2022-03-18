# case_watcher
Intended to make it easier to see potential new cases being opened via a keyword search

# Configuration File Entries
* `ses_from_email`: This is the 'from' email address to use when sending the email report, it needs to be verified with the Amazon SES service.

# Credentials
## Case Repository
URL, Username, and Password are needed for the endpoint giving us case information

## Google Cloud IAM account
Through the configuration file we expect to be given a google iam account and a spreadsheet link to update

## Amazon IAM Account
We need AWS IAM credentials in the environment so we may email a report via Amazon SES
