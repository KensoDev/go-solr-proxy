git tag $1
gox -os="linux" -tags="$1" -output="out"