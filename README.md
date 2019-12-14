export YOUR_INPUT_BUCKET_NAME=tsukada-input
gsutil mb gs://$YOUR_INPUT_BUCKET_NAME

export YOUR_OUTPUT_BUCKET_NAME=tsukada-output
gsutil mb gs://$YOUR_OUTPUT_BUCKET_NAME


mkdir project
cd project
git clone https://github.com/GoogleCloudPlatform/golang-samples.git
cd golang-samples/functions/imagemagick/

gcloud functions deploy BlurOffensiveImages --runtime go111 --trigger-bucket $YOUR_INPUT_BUCKET_NAME --set-env-vars BLURRED_BUCKET_NAME=$YOUR_OUTPUT_BUCKET_NAME

gsutil cp *.jpg gs://$YOUR_INPUT_BUCKET_NAME

