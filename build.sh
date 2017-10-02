gomobile bind -target=ios github.com/mebusy/godict/dictUtils

DEST_PATH="../../../../../dict2017/"
rm -rf ${DEST_PATH}*.framework
mv -f *.framework ${DEST_PATH}


