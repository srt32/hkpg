# hkpg

Automated Heroku Postgres backups written in go

Inspired by https://github.com/kbaum/heroku-database-backups

Deploy with ideas in https://stackoverflow.com/a/43265218/1949363

Manage deps / deploys with: https://devcenter.heroku.com/articles/deploying-go


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
                "arn:aws:s3:::bucket-name/*"
            ]
        }
    ]
}
```
