go install omo.msa.user
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.user _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.user.tar.gz ./*
mv msa.user.tar.gz ../
cd ../
rm -rf _build
