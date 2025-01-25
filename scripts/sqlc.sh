# Find all the folders that contain sqlc.yaml
find . -name sqlc.yaml

# Run sqlc generate in each folder
for dir in $(find . -name sqlc.yaml); do
  echo "Running sqlc generate in $dir"
  cd $(dirname $dir)
  sqlc generate
  cd -
done
