# hkpg
### Automated Archiving of Heroku Postgres Backups Written in Go

This package contains a single command that can be run manually or via the
Heroku Scheduler to pull down the most recent Heroku Postgres backup from
another Heroku application and upload it to a specified S3 Bucket.

This work is inspired by [similar
work](https://github.com/kbaum/heroku-database-backups) done via bash.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Configuring Heroku

* Create a new application on Heroku via `heroku create`
* Clone the repo via `git clone git@github.com:srt32/hkpg.git`
* Deploy the repo via `git push heroku master`
* Set up environment variables as described below
* Setup Heroku Scheduler to run `bin/backup` on whichever interval you choose

## Environment Variables

The below environment variables are required unless otherwise specified:

* HKPG_HEROKU_AUTH_TOKEN="copy-from-heroku-dashboard"
* HKPG_HEROKU_APP_NAME="my-app-name"
* S3_BUCKET_NAME="your-bucket-name"
* AWS_ACCESS_KEY_ID="key-id"
* AWS_SECRET_KEY="secret-id"
* AWS_REGION="us-east-1" (optional, defaults to us-west-1)

## Configuring S3 Bucket Permissions

You will need an IAM user with "s3:PutObject" and "s3:PutObjectAcl" permissions
on the root:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1504540641000",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:PutObjectAcl"
            ],
            "Resource": [
                "arn:aws:s3:::your-bucket-name/*"
            ]
        }
    ]
}
```
