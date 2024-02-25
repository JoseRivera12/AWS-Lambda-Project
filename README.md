# Stori Tech Challenge

AWS Lambda function built using the AWS Cloud Development Kit (CDK). The function is designed to read CSV files containing a user's transaction data. Once parsed, it sends an email to the user, summarizing their monthly transactions and providing their current account balance. Additionally, the function includes functionality to store these transactions in a database for future reference.

<img src="https://github.com/JoseRivera12/AWS-Lambda-Project/assets/23196720/97b6845e-b2a1-4da7-814f-708dfc7014a2" width=60% height=60%>

## Installation

- Setup CDK
[AWS CDK](https://docs.aws.amazon.com/cdk/v2/guide/getting_started.html)
- Setup ad validate an email account in SES
[AWS SES](https://docs.aws.amazon.com/ses/latest/dg/setting-up.html)
- Setup a relational postgresql database in AWS
[AWS RDS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html)

In the file **stori_project.go** assign the env variables.
```
    "DATABASE_NAME":     jsii.String("db_name"),
    "DATABASE_USER":     jsii.String("db_user"),
    "DATABASE_HOST":     jsii.String("db_host"),
    "DATABASE_PORT":     jsii.String("db_port"),
    "DATABASE_PASSWORD": jsii.String("db_password"),
    "SES_EMAIL":         jsii.String("db_email"),
```

Deploy the stack in AWS 
```
cdk bootstrap
cdk deploy
```

***By default the customer information and email template are readed from S3, you can find the files in the root of the project.**

If you need to setup another user information, you need to upload a txt file **user.txt** with the following format. **If there's not a user.txt file by default the lambda function read user information from the default S3.**

```
Customer Name
customer@test.com
```

## Test lambda function

To test the function it is necessary to upload one or multiple csv files to the bucket that was created next to the stack. You can find a test file in the root directory of the project **user_transactions.csv**.

## Running Tests

Run test suit inside the lambda-handler for lambda function

```
go test ./... 
```

Run unit tests in root project for CDK tests
```
go test ./... 
```

## Useful commands

Useful AWS commands
```
cdk diff // compare deployed stack with current state
cdk synth // emits the synthesized CloudFormation template
```

You can generate a .csv file for testing with the following command
```
// python csv_generator.py file_name num_transactions
python csv_generator.py user_transactions 20000
```

## CSV File Format
```
Id,Date,Transaction,Currency,Description
1,2024-10-14T11:54:41,+594.83,USD,Salary
```

## Example Email Sended
![example_email](https://github.com/JoseRivera12/AWS-Lambda-Project/assets/23196720/0a18352a-5f4b-4a28-a949-19f9ae4790b7)