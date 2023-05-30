# Sebuung Competition Web

When you connect to port 80, a web page for file submission is displayed. When a file is submitted, the golang application uploads the file to the bucket set in the Configfile. Additionally, the file is stored locally in the ./files/ directory. Since the application uses the default AWS credentials, please set up a default profile that can access S3 using the aws configure command.

# Configuration

Save the S3 Bucket Name where the file will be stored in config.json of repo.

```jsx
{
    "BUCKET_NAME" : "your-bucket-name" 
}
```

# Build & Deploy Guide

After receiving the repository with git clone, execute it as follows.

```jsx
go build .
chmod +x Sebuung-Competition-Web
./Sebuung-Competition-Web
```
