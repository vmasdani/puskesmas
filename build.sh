rm -rf dist &&\ 
mkdir -p dist/{landing,admin,news} &&\
cd landing &&\
npm run build &&\
cp -r build/* ../dist/landing &&\
cd ../admin &&\
npm run build &&\
cp -r build/* ../dist/admin &&\
cd .. &&\
go build &&\
mv puskesmas dist