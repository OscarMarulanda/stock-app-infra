resource "random_id" "suffix" {
  byte_length = 4
}

resource "aws_s3_bucket" "frontend" {
  bucket = "stocks-bucket-${random_id.suffix.hex}"  # Replace with your desired bucket name

  # Remove the `acl` attribute; AWS recommends using bucket policies for access control.
}

resource "aws_s3_bucket_website_configuration" "frontend_website" {
  bucket = aws_s3_bucket.frontend.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}