cd auth_service && go mod tidy && go mod vendor && go vet ./... && go mod tidy && cd ../
echo "💚 auth_service done"
cd notification_service && go mod tidy && go mod vendor && go vet ./... && go mod tidy && cd ../
echo "💚 notification_service done"
cd rest_user_service && go mod tidy && go mod vendor && go vet ./... && go mod tidy && cd ../
echo "💚 rest_user_service done"
cd system_service && go mod tidy && go mod vendor && go vet ./... && go mod tidy && cd ../
echo "💚 system_service done"
echo "🩵 all done"