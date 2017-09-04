# hkpg
### Automated archiving of Heroku Postgres backups written in Go

This package contains a single command that can be run manually or via the
Heroku Scheduler to pull down the most recent Heroku Postgres backup from
another Heroku application and upload it to a specified S3 Bucket.

This work is inspired by [similar
work](https://github.com/kbaum/heroku-database-backups) done via bash.

## Environment Variables

The below environment variables are required unless otherwise specified:

* HKPG_HEROKU_AUTH_TOKEN="copy-from-heroku-dashboard"
* HKPG_HEROKU_APP_NAME="my-app-name"
* S3_BUCKET_NAME="your-bucket-name"
* AWS_ACCESS_KEY_ID="key-id"
* AWS_SECRET_KEY="secret-id"
* AWS_REGION="us-east-1" (optional, defaults to us-west-1)

## Configuring Heroku Scheduler

Setup Heroku Scheduler to run `bin/backup`.

## Configuring S3 bucket permissions

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
