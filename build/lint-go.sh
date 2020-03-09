if [ -z "$(which golint)" ]
then
      go get -u golang.org/x/lint/golint
fi


golint ../...