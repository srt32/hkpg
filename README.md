# hkpg

Automated archiving of Heroku Postgres backups written in Go

Inspired by https://github.com/kbaum/heroku-database-backups

Deploy with ideas in https://stackoverflow.com/a/43265218/1949363

Manage deps / deploys with: https://devcenter.heroku.com/articles/deploying-go

TODO:
* wrap log statements in HKPG_DEBUG check
* add a license
* add tests around core behavior

## Environment Variables

HEROKU_AUTH_TOKEN=""
HEROKU_APP_NAME=""
AWS_ACCESS_KEY_ID=""
AWS_SECRET_KEY=""
S3_BUCKET_NAME="" (optional, defaults to us-west-1)

## Configuring Heroku Scheduler

TODO

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
